package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
)

func AuthenticateHandler(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	err := r.ParseForm()

	if err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Error parsing form data")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if password == "" || username == "" {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Missing sign in info")
		return
	}

	usernameLength := utf8.RuneCountInString(username)

	if usernameLength < 3 || usernameLength > 20 {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Username length requirements not met")
		return
	}

	passwordLength := utf8.RuneCountInString(password)

	if passwordLength < 8 || passwordLength > 64 {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Password length requirements not met")
		return
	}

	row := dbContext.Connection.QueryRow(
		context.Background(),
		"SELECT account_id, password_hash, session_id FROM Accounts WHERE account_name=$1",
		username)

	var accountId int
	var hash string
	var sessionId pgtype.Text

	scanErr := row.Scan(&accountId, &hash, &sessionId)

	if scanErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Querying account info failed")
		return
	}

	match, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Compare hash failed")
		return
	}

	if !match {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Password hashes do not match")
		return
	}

	newSessionId, err := helpers.GenerateBase64RandomId(48)

	if err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to generate session id")
		return
	}

	expirationDate := time.Now().AddDate(0, 0, 1)

	insertNewSession := helpers.ParameterizedQuery{
		Sql:    "INSERT INTO Sessions (session_id, account_id, expires_at) VALUES ($1, $2, $3)",
		Params: []any{newSessionId, accountId, expirationDate},
	}

	updateAccount := helpers.ParameterizedQuery{
		Sql:    "UPDATE Accounts SET session_id = $1 WHERE account_id = $2",
		Params: []any{newSessionId, accountId},
	}

	deleteOldSession := helpers.ParameterizedQuery{
		Sql:    "DELETE FROM Sessions WHERE session_id = $1",
		Params: []any{sessionId},
	}

	transactionErr := helpers.RunTransaction(dbContext, []*helpers.ParameterizedQuery{
		&insertNewSession,
		&updateAccount,
		&deleteOldSession,
	})

	if transactionErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, "Unable to commit session transaction")
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    newSessionId,
		Path:     "/",
		HttpOnly: true,
		// Secure
		// SameSite
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/")
}
