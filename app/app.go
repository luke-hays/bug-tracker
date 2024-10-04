package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Incredibly messy logic to just to test connecting to a db
	dbURL := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}
	fmt.Println("Connected to the database successfully!")

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
