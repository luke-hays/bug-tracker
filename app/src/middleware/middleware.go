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

			if err == nil && sessionCookie != nil {
				// TODO Check Expiration Date
				session := dbContext.Connection.QueryRow(
					context.Background(),
					"SELECT account_id FROM Sessions WHERE session_id = $1;",
					sessionCookie.Value,
				)

				var accountId uint64

				scanErr := session.Scan(&accountId)

				if scanErr != nil {
					fmt.Printf("%v\n", scanErr)
					http.Redirect(w, r, "/signin", http.StatusSeeOther)
					return
				}

				key := helpers.AccountIdKey("account_id")
				ctx := context.WithValue(r.Context(), key, accountId)
				r = r.WithContext(ctx)
			} else {
				fmt.Printf("%v\n", err)
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
