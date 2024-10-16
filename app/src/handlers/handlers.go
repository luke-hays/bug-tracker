package handlers

import (
	"app/src/db"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"unicode/utf8"

	"github.com/alexedwards/argon2id"
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
		"SELECT password_hash FROM Accounts WHERE account_name=$1",
		username)

	var hash string

	scanError := row.Scan(&hash)

	if scanError != nil {
		fmt.Println("No results found")
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

	w.WriteHeader(http.StatusOK)
}
