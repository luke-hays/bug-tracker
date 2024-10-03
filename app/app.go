package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Temporary solution just to test composing a base layout, a page, and multiple components together
	tmpl := template.Must(template.ParseFiles("templates/layouts/base.html", "templates/pages/home.html", "components/example-btn.html"))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	// Need to strip the static prefix from the path so that we ca serve static assets
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":8080", router)
}
