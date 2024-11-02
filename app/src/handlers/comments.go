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

type comments struct {
	CommentId   uint64    `db:"comment_id"`
	BugId       uint64    `db:"bug_id"`
	Author      uint64    `db:"author"`
	CommentDate time.Time `db:"comment_date"`
	Comment     string    `db:"comment"`
}

func CreateComment(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	parseFormErr := r.ParseForm()

	if parseFormErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, parseFormErr.Error())
		return
	}

	// TODO make this a constant or maybe helper
	key := helpers.AccountIdKey("account_id")
	accountId := r.Context().Value(key)

	bugId := r.FormValue("bug_id")
	author := accountId
	commentDate := time.Now()
	comment := r.FormValue("comment")

	bugIdNum, atoiConvErr := strconv.Atoi(bugId)

	if atoiConvErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, atoiConvErr.Error())
		return
	}

	insertComment := helpers.ParameterizedQuery{
		Sql:    "INSERT INTO Comments (bug_id, author, comment_date, comment) VALUES ($1, $2, $3, $4);",
		Params: []any{bugIdNum, author, commentDate, comment},
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

	tmpl := template.Must(template.ParseFiles("src/components/bug-comments.html"))
	tmpl.Execute(w, bugComments)
}
