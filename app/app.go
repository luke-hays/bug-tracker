package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"app/src/db"
	"app/src/routes"
)

func main() {
	dbContext, error := db.Init()

	if error != nil {
		log.Fatalf("Unable to connect to database: %v\n", error)
	}

	router := mux.NewRouter()
	routes.RegisterRoutes(router, dbContext)

	// Need to strip the static prefix from the path so that we can serve static assets
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":8080", router)
}
