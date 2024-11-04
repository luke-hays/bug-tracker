package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type bug struct {
	BugId        uint64    `db:"bug_id"`
	DateReported time.Time `db:"date_reported"`
	Summary      *string   `db:"summary"`
	Description  *string   `db:"description"`
	Resolution   *string   `db:"resolution"`
	ReportedBy   uint64    `db:"reported_by"`
	AssignedTo   *uint64   `db:"assigned_to"`
	VerifiedBy   *uint64   `db:"verified_by"`
	Status       string    `db:"status"`
	Priority     *string   `db:"priority"`
	Hours        *float64  `db:"hours"`
}

func CreateBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	// TODO remove this check along with the account id retrieval
	// Middleware can handle this now
	sessionId, err := r.Cookie("session")

	if err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusUnauthorized, "Unable to authorize session id")
		return
	}

	parseFormErr := r.ParseForm()

	if parseFormErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Error parsing form data")
		return
	}

	dateReported := time.Now()
	summary := r.FormValue("summary")                          // Optional
	description := r.FormValue("description")                  // Optional
	assignedTo := helpers.ParseInt(r.FormValue("assigned_to")) // Optional
	status := "NEW"
	priority := r.FormValue("priority")                 // Optional
	hours := helpers.ParseFloat64(r.FormValue("hours")) // Optional

	session := dbContext.Connection.QueryRow(
		context.Background(),
		"SELECT account_id FROM Sessions WHERE session_id = $1;",
		sessionId.Value,
	)

	var reportedBy int

	scanErr := session.Scan(&reportedBy)

	if scanErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusUnauthorized, "Unable to create bug - cannot find account id")
		return
	}

	insertNewBug := helpers.ParameterizedQuery{
		Sql:    "INSERT INTO Bugs (date_reported, summary, description, reported_by, assigned_to, status, priority, hours) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
		Params: []any{dateReported, summary, description, reportedBy, assignedTo, status, priority, hours},
	}

	// TODO It would be helpful to return things from this like ids
	transactionErr := helpers.RunTransaction(dbContext, []*helpers.ParameterizedQuery{&insertNewBug})

	if transactionErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to commit transaction for inserting a new bug")
		return
	}

	// TODO Add method to send select menu for users
	w.WriteHeader(http.StatusCreated)
	tmpl := template.Must(template.ParseFiles("src/components/bug-info.html"))
	tmpl.ExecuteTemplate(w, "bug-info", &bug{})
}

func UpdateBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	parseFormErr := r.ParseForm()

	if parseFormErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, parseFormErr.Error())
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
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, transactionErr.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetBugs(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	// TODO implement pagination
	// TODO Time isn't recorded
	bugRecords, queryBugRecordsErr := dbContext.Connection.Query(context.Background(), "SELECT * FROM Bugs ORDER BY date_reported DESC, bug_id DESC")

	if queryBugRecordsErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, queryBugRecordsErr.Error())
		return
	}

	bugs, collectRecordsErr := pgx.CollectRows(bugRecords, pgx.RowToStructByName[bug])

	if collectRecordsErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, collectRecordsErr.Error())
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"src/templates/layouts/base.html",
		"src/templates/pages/bugs.html",
		"src/components/bug-info.html",
	))

	tmpl.Execute(w, bugs)
}

func GetBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext, bugId string) {
	rows, queryBugErr := dbContext.Connection.Query(context.Background(), "SELECT * FROM Bugs WHERE bug_id = $1", bugId)

	if queryBugErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusNotFound, queryBugErr.Error())
		return
	}

	bugRecord, collectRecordErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[bug])

	if collectRecordErr != nil {
		// TODO this could technically be a not found, or duplicate id issue
		helpers.WriteAndLogHeaderStatus(w, http.StatusNotFound, collectRecordErr.Error())
		return
	}

	encodedBug, encodingErr := json.Marshal(bugRecord)

	if encodingErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, encodingErr.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedBug)
}
