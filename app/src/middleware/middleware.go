package middleware

import (
	helpers "app/src/utils"
	"fmt"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if helpers.RequestHasValidSession(r) {
			// If the cookie is missing, redirect to /signin
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
