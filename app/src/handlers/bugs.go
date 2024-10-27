package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type bug struct {
	BugId        uint32    `db:"bug_id"`
	DateReported time.Time `db:"date_reported"`
	Summary      *string   `db:"summary"`
	Description  *string   `db:"description"`
	Resolution   *string   `db:"resolution"`
	ReportedBy   uint32    `db:"reported_by"`
	AssignedTo   *uint32   `db:"assigned_to"`
	VerifiedBy   *uint32   `db:"verified_by"`
	Status       string    `db:"status"`
	Priority     *string   `db:"priority"`
	Hours        *float64  `db:"hours"`
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
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Error querying bugs")
		return
	}

	bugs, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[bug])

	if collectErr != nil {
		fmt.Printf("%s\n", collectErr.Error())
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to collect rows")
		return
	}

	tmpl := template.Must(template.ParseFiles("src/templates/layouts/base.html", "src/templates/pages/bugs.html"))
	tmpl.Execute(w, bugs)
}

func GetBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext, bugId string) {
	// Get specific bug by id
	fmt.Printf("%s\n", bugId)
	rows, queryErr := dbContext.Connection.Query(context.Background(), "SELECT * FROM Bugs WHERE bug_id = $1", bugId)

	if queryErr != nil {
		fmt.Printf("%s\n", queryErr.Error())
		helpers.WriteAndLogHeaderStatus(w, http.StatusNotFound, "Unable to query for bug")
		return
	}

	bugRecord, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[bug])

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		helpers.WriteAndLogHeaderStatus(w, http.StatusNotFound, "Unable to find bug")
		return
	}

	encodedBug, encodingErr := json.Marshal(bugRecord)

	if encodingErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to encode bugs")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedBug)
}
