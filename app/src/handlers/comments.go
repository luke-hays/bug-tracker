package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type oneComment struct {
	Author      uint64
	CommentDate time.Time
	Comment     string
}

type comments struct {
	CommentId   uint64    `db:"comment_id"`
	BugId       uint64    `db:"bug_id"`
	Author      uint64    `db:"author"`
	CommentDate time.Time `db:"comment_date"`
	Comment     string    `db:"comment"`
}

type bugCommentsPair struct {
	BugId    string
	Comments []comments
}

func CreateComment(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext, bugId string) {
	parseFormErr := r.ParseForm()

	if parseFormErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, parseFormErr.Error())
		return
	}

	// TODO make this a constant or maybe helper
	key := helpers.AccountIdKey("account_id")
	accountId := r.Context().Value(key)

	authorId, ok := accountId.(uint64)

	if !ok {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Author ID is not an valid ID")
		return
	}

	commentDate := time.Now()
	comment := r.FormValue("comment")

	bugIdNum, atoiConvErr := strconv.Atoi(bugId)

	if atoiConvErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, atoiConvErr.Error())
		return
	}

	insertComment := helpers.ParameterizedQuery{
		Sql:    "INSERT INTO Comments (bug_id, author, comment_date, comment) VALUES ($1, $2, $3, $4);",
		Params: []any{bugIdNum, authorId, commentDate, comment},
	}

	transactionErr := helpers.RunTransaction(
		dbContext,
		[]*helpers.ParameterizedQuery{&insertComment},
	)

	if transactionErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, transactionErr.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	tmpl := template.Must(template.ParseFiles("src/components/bug-comment.html"))
	tmpl.ExecuteTemplate(w, "bug-comment", &oneComment{Author: authorId, CommentDate: commentDate, Comment: comment})
}

func GetComments(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext, bugId string) {
	rows, queryErr := dbContext.Connection.Query(context.Background(), "SELECT * from Comments WHERE bug_id = $1", bugId)

	if queryErr != nil {
		fmt.Printf("%s\n", queryErr.Error())
		helpers.WriteAndLogHeaderStatus(w, http.StatusNotFound, queryErr.Error())
		return
	}

	bugComments, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[comments])

	if collectErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, collectErr.Error())
		return
	}

	tmpl := template.Must(template.ParseFiles("src/components/bug-comments.html", "src/components/bug-comment.html"))
	tmpl.Execute(w, &bugCommentsPair{BugId: bugId, Comments: bugComments})
}
