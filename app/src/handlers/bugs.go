package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"net/http"
	"time"
)

func CreateBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext, sessionId string) {
	err := r.ParseForm()

	if err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Error parsing form data")
		return
	}

	dateReported := time.Now()
	summary := r.FormValue("summary")
	description := r.FormValue("description")
	assignedTo := r.FormValue("assigned_to")
	status := "NEW"
	priority := r.FormValue("priority")
	hours := r.FormValue("hours")

	session := dbContext.Connection.QueryRow(
		context.Background(),
		"SELECT account_id FROM Sessions WHERE session_id = $1;",
		sessionId,
	)

	var reportedBy int

	scanErr := session.Scan(&reportedBy)

	if scanErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to create bug - cannot find account id")
		return
	}

	insertNewBug := helpers.ParameterizedQuery{
		Sql:    "INSERT INTO Bugs (date_reported, summary, description, resolution, reported_by, assigned_to, status, priority, hours) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
		Params: []any{dateReported, summary, description, reportedBy, assignedTo, status, priority, hours},
	}

	transactionErr := helpers.RunTransaction(dbContext, []*helpers.ParameterizedQuery{&insertNewBug})

	if transactionErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to commit transaction for inserting a new bug")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	// Get bug id and session id
	// Can update
	// description
	// resolution
	// summary
	// assigned to
	// verified by
	// status
	// priority
	// hours
}

func GetBugs(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	// Get all bugs, use pagination with limit of 10. Offset is just pageNum * limit
}

func GetBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	// Get specific bug by id
}
