package routes

import (
	"app/src/db"
	"app/src/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, dbContext *db.DatabaseContext) {
	router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		handlers.SignInHandler(w, r)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, dbContext)
	}).Methods("GET")
}
