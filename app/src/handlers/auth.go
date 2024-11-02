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
	parseErr := r.ParseForm()

	if parseErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, parseErr.Error())
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if password == "" || username == "" {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Missing authentication info")
		return
	}

	usernameLength := utf8.RuneCountInString(username)
	minUsernameLength := 3
	maxUsernameLength := 20

	if usernameLength < minUsernameLength || usernameLength > maxUsernameLength {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Username length invalid")
		return
	}

	passwordLength := utf8.RuneCountInString(password)
	minPasswordLength := 8
	maxPasswordLength := 64

	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Password length invalid")
		return
	}

	accountRecord := dbContext.Connection.QueryRow(
		context.Background(),
		"SELECT account_id, password_hash, session_id FROM Accounts WHERE account_name=$1",
		username)

	var accountId int
	var hash string
	var sessionId pgtype.Text

	scanAccountRecordErr := accountRecord.Scan(&accountId, &hash, &sessionId)

	if scanAccountRecordErr != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusNotFound, scanAccountRecordErr.Error())
		return
	}

	passwordMatch, hashCompareError := argon2id.ComparePasswordAndHash(password, hash)

	if hashCompareError != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, hashCompareError.Error())
		return
	}

	if !passwordMatch {
		helpers.WriteAndLogHeaderStatus(w, http.StatusBadRequest, "Invalid password")
		return
	}

	newSessionId, generateBase64Err := helpers.GenerateBase64RandomId(48)

	if generateBase64Err != nil {
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, generateBase64Err.Error())
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
		helpers.WriteAndLogHeaderStatus(w, http.StatusInternalServerError, transactionErr.Error())
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
	w.Header().Set("Location", "/bugs")
}
