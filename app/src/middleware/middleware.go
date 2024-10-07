package middleware

import (
	"fmt"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling middleware logger")
		next.ServeHTTP(w, r)
	})
}
