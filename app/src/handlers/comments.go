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
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	key := helpers.AccountIdKey("account_id")
	accountId := r.Context().Value(key)

	fmt.Printf("accountId - %d\n", accountId)

	bugId := r.FormValue("bug_id")
	author := accountId
	commentDate := time.Now()
	comment := r.FormValue("comment")

	bugIdNum, convErr := strconv.Atoi(bugId)

	if convErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Unable to parse bug id")
		return
	}

	fmt.Printf("%d, %d, %s, %s\n", bugIdNum, author, commentDate, comment)

	insertComment := helpers.ParameterizedQuery{
		Sql:    "INSERT INTO Comments (bug_id, author, comment_date, comment) VALUES ($1, $2, $3, $4);",
		Params: []any{bugIdNum, author, commentDate, comment},
	}

	transactionErr := helpers.RunTransaction(
		dbContext,
		[]*helpers.ParameterizedQuery{&insertComment},
	)

	if transactionErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to commit transaction for inserting a new comment")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetComments(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext, bugId string) {
	rows, queryErr := dbContext.Connection.Query(context.Background(), "SELECT * from Comments WHERE bug_id = $1", bugId)

	if queryErr != nil {
		fmt.Printf("%s\n", queryErr.Error())
		helpers.WriteAndLogHeaderStatus(w, http.StatusNotFound, "Unable to query for comments")
		return
	}

	bugComments, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[comments])

	if collectErr != nil {
		fmt.Printf("%s\n", collectErr.Error())
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to collect rows")
		return
	}

	tmpl := template.Must(template.ParseFiles("src/components/bug-comments.html"))
	tmpl.Execute(w, bugComments)
}
