package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type bug struct {
	Bug_id        uint32
	Date_reported time.Time
	Summary       sql.NullString
	Description   sql.NullString
	Resolution    sql.NullString
	Reported_by   uint32
	Assigned_to   sql.NullInt32
	Verified_by   sql.NullInt32
	Status        string
	Priority      sql.NullString
	Hours         sql.NullFloat64
}

func CreateBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	sessionId, err := r.Cookie("session")

	if err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to authorize session id")
		return
	}

	parseFormErr := r.ParseForm()

	if parseFormErr != nil {
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
		helpers.WriteAndLogHeaderStatus(w, http.StatusUnauthorized, "Unable to create bug - cannot find account id")
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
	err := r.ParseForm()

	if err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Error parsing form data")
		return
	}

	bugId := r.FormValue("bug_id")
	summary := r.FormValue("summary")
	description := r.FormValue("description")
	resolution := r.FormValue("resolution")
	assignedTo := r.FormValue("assigned_to")
	verifiedBy := r.FormValue("verified_by")
	status := r.FormValue("status")
	priority := r.FormValue("priority")
	hours := r.FormValue("hours")

	updateBug := helpers.ParameterizedQuery{
		Sql:    "UPDATE Bugs SET (summary, description, resolution, assigned_to, verified_by, status, priority, hours) = ($2, $3, $4, $5, $6, $7, $8, $9) WHERE bug_id = $1",
		Params: []any{bugId, summary, description, resolution, assignedTo, verifiedBy, status, priority, hours},
	}

	transactionErr := helpers.RunTransaction(dbContext, []*helpers.ParameterizedQuery{&updateBug})

	if transactionErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Error updating bug")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetBugs(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	// For now encode as json and send it back
	// also probably want to paginate
	rows, err := dbContext.Connection.Query(context.Background(), "SELECT * FROM Bugs")

	if err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Error updating bug")
		return
	}

	bugs, collectErr := pgx.CollectRows(rows, pgx.RowTo[bug])

	if collectErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to collect rows")
		return
	}

	encodedBugs, encodingErr := json.Marshal(bugs)

	if encodingErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to encode bugs")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedBugs)
}

func GetBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	// Get specific bug by id
}
