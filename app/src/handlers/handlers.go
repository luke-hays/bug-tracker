package handlers

import (
	"app/src/db"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, dbContext *db.DatabaseContext) {
	tmpl := template.Must(template.ParseFiles("src/templates/layouts/base.html", "src/templates/pages/home.html", "src/components/example-btn.html"))
	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("src/templates/layouts/base.html", "src/templates/pages/home.html", "src/components/example-btn.html"))
	tmpl.Execute(w, nil)
}
