package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"app/src/db"
	"app/src/middleware"
	"app/src/routes"
)

func main() {
	dbContext, error := db.Init()

	if error != nil {
		log.Fatalf("Unable to connect to database: %v\n", error)
	}

	pingErr := dbContext.Connection.Ping(context.Background())

	if pingErr != nil {
		log.Fatalf("Pinging database failed: %v\n", error)
	}

	router := mux.NewRouter()

	routes.RegisterRoutes(router, dbContext)

	router.Use(middleware.Authenticator)
	router.Use(middleware.Logger)

	// Need to strip the static prefix from the path so that we ca serve static assets
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":8080", router)
}
