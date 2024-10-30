package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

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

	_, insertErr := dbContext.Connection.Query(
		context.Background(),
		"INSERT INTO Comments (bug_id, author, comment_date, comment) VALUES ($1, $2, $3, $4)",
		bugIdNum, author, commentDate, comment,
	)

	if insertErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to insert new comment")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
