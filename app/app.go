package main

import (
	"html/template"
	"net/http"
)

func main() {
	// Temporary solution just to test composing a base layout, a page, and multiple components together
	tmpl := template.Must(template.ParseFiles("templates/layouts/base.html", "templates/pages/home.html", "components/example-btn.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.ListenAndServe(":8080", nil)
}
