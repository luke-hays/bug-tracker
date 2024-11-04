package routes

import (
	"app/src/db"
	"app/src/handlers"
	"app/src/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func registerSecureRoutes(router *mux.Router, dbContext *db.DatabaseContext) {
	secureRoutes := router.NewRoute().Subrouter()

	secureRoutes.Use(middleware.Authenticator(dbContext))

	secureRoutes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO Redirect to Bugs
		handlers.HomeHandler(w, r, dbContext)
	}).Methods("GET")

	/***** Bug Routes *****/
	secureRoutes.HandleFunc("/bugs", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetBugs(w, r, dbContext)
	}).Methods("GET")

	secureRoutes.HandleFunc("/bugs", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateBug(w, r, dbContext)
	}).Methods("POST")

	secureRoutes.HandleFunc("/bugs/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bugId := vars["id"]
		handlers.GetBug(w, r, dbContext, bugId)
	}).Methods("GET")

	secureRoutes.HandleFunc("/bugs/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateBug(w, r, dbContext)
	}).Methods("PUT")

	/***** Comment Routes *****/
	secureRoutes.HandleFunc("/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bugId := vars["id"]
		handlers.CreateComment(w, r, dbContext, bugId)
	}).Methods("POST")

	secureRoutes.HandleFunc("/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bugId := vars["id"]
		handlers.GetComments(w, r, dbContext, bugId)
	}).Methods("GET")
}

func registerPublicRoutes(router *mux.Router, dbContext *db.DatabaseContext) {
	router.HandleFunc("/api/authenticate", func(w http.ResponseWriter, r *http.Request) {
		handlers.AuthenticateHandler(w, r, dbContext)
	}).Methods("POST")

	// TODO Can we display a timeout message? Should we?
	router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		handlers.SignInHandler(w, r)
	}).Methods("GET")
}

func RegisterRoutes(router *mux.Router, dbContext *db.DatabaseContext) {
	// Middleware used by all routes
	router.Use(middleware.Logger)

	registerSecureRoutes(router, dbContext)
	registerPublicRoutes(router, dbContext)
}
