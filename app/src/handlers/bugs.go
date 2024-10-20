package handlers

import (
	"app/src/db"
	"net/http"
)

func CreateBug(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	// Get account id from session id
	// Date reported is current timestampe
	// Reported by is an existing account id
	// Bug should start on new as default - need to add statuses in bug status
	// Optionally start with
	// summary
	// description
	// assigned to
	// priority
	// hours
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
