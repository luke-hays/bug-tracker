package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"fmt"
	"net/http"
)

func CreateComment(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	parseFormErr := r.ParseForm()

	if parseFormErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	key := helpers.AccountIdKey("account_id")
	accountId := r.Context().Value(key)

	fmt.Println(accountId)

	// bugId := helpers.ParseInt(r.FormValue("bug_id"))
	// author := helpers.ParseInt(r.FormValue("author"))
	// commentDate := time.Now()
	// comment := r.FormValue("comment")

}
