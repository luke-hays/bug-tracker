package handlers

import (
	"app/src/db"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("src/templates/layouts/base.html", "src/templates/pages/home.html", "src/components/example-btn.html"))
	tmpl.Execute(w, nil)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("src/templates/layouts/base.html", "src/templates/pages/signin.html"))
	tmpl.Execute(w, nil)
}

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

	scanError := row.Scan(&accountId, &hash, &sessionId)

	if scanError != nil {
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

	// 48 bytes = 64 base64 characters
	randomBytes := make([]byte, 48)
	_, sessionIdError := rand.Read(randomBytes)

	if sessionIdError != nil {
		fmt.Println("Unable to generate session id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newSessionId := base64.URLEncoding.EncodeToString(randomBytes)
	expirationDate := time.Now().AddDate(0, 0, 1)

	tx, err := dbContext.Connection.Begin((context.Background()))

	if err != nil {
		fmt.Println("Unable to begin session creation transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("session id length: %d\n", utf8.RuneCountInString(newSessionId))

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
		Name:  "session",
		Value: newSessionId,
		Path:  "/",

		// Secure
		HttpOnly: true,
		// SameSite
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/")
}
