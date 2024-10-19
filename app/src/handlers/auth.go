package handlers

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
)

func AuthenticateHandler(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if password == "" || username == "" {
		fmt.Println("Missing sign in info")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usernameLength := utf8.RuneCountInString(username)

	if usernameLength < 3 || usernameLength > 20 {
		fmt.Println("Username length requirements not met")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passwordLength := utf8.RuneCountInString(password)

	if passwordLength < 8 || passwordLength > 64 {
		fmt.Println("Password length requirements not met")
		w.WriteHeader(http.StatusBadRequest)
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
		fmt.Println("Querying account info failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		fmt.Println("Unable to hash potential password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !match {
		fmt.Println("Password hashes do not match")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newSessionId, err := helpers.GenerateBase64RandomId(48)

	if err != nil {
		fmt.Println("Unable to generate session id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expirationDate := time.Now().AddDate(0, 0, 1)

	tx, err := dbContext.Connection.Begin((context.Background()))

	if err != nil {
		fmt.Println("Unable to begin session creation transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(context.Background())

	tx.Exec(
		context.Background(),
		"INSERT INTO Sessions (session_id, account_id, expires_at) VALUES ($1, $2, $3)",
		newSessionId, accountId, expirationDate)

	tx.Exec(
		context.Background(),
		"UPDATE Accounts SET session_id = $1 WHERE account_id = $2",
		newSessionId, accountId)

	tx.Exec(
		context.Background(),
		"DELETE FROM Sessions WHERE session_id = $1",
		sessionId)

	err = tx.Commit(context.Background())

	if err != nil {
		fmt.Println("Unable to commit session transaction")
		w.WriteHeader(http.StatusInternalServerError)
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
