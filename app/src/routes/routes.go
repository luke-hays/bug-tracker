package routes

import (
	"app/src/db"
	"app/src/handlers"
	"app/src/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, dbContext *db.DatabaseContext) {
	router.HandleFunc("/api/authenticate", func(w http.ResponseWriter, r *http.Request) {
		handlers.AuthenticateHandler(w, r, dbContext)
	}).Methods("POST")

	router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		handlers.SignInHandler(w, r)
	}).Methods("GET")

	sr := router.NewRoute().Subrouter()

	sr.Use(middleware.Authenticator)

	sr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r)
	}).Methods("GET")
}
