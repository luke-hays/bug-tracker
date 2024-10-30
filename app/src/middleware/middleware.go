package middleware

import (
	"app/src/db"
	helpers "app/src/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func Authenticator(dbContext *db.DatabaseContext) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie("session")

			// fmt.Printf("session cookie %s\n", sessionCookie.Value)

			if err == nil && sessionCookie != nil {
				// Need to check that the session exists and is not expired
				session := dbContext.Connection.QueryRow(
					context.Background(),
					"SELECT account_id FROM Sessions WHERE session_id = $1;",
					sessionCookie.Value,
				)

				var accountId uint32

				scanErr := session.Scan(&accountId)

				if scanErr != nil {
					fmt.Printf("%s\n", scanErr.Error())
					http.Redirect(w, r, "/signin", http.StatusSeeOther)
					return
				}

				key := helpers.AccountIdKey("account_id")
				ctx := context.WithValue(r.Context(), key, accountId)
				r = r.WithContext(ctx)
			} else {
				fmt.Printf("%s\n", err.Error())
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
