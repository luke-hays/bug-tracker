package helpers

import (
	"app/src/db"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
)

type ParameterizedQuery struct {
	Sql    string
	Params []any
}

// Check for the session cookie
func RequestHasValidSession(r *http.Request) bool {
	sessionCookie, err := r.Cookie("session")
	return err != nil && sessionCookie != nil
}

// Generate a random id for a specified number of bytes.
// eg 48 bytes == 64 char string
func GenerateBase64RandomId(byteNum int) (string, error) {
	randomBytes := make([]byte, byteNum)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func WriteAndLogHeaderStatus(w http.ResponseWriter, status int, message string) {
	fmt.Println(message)
	w.WriteHeader(status)
}

func RunTransaction(dbContext *db.DatabaseContext, queries []*ParameterizedQuery) error {
	if len(queries) == 0 {
		return nil
	}

	tx, err := dbContext.Connection.Begin((context.Background()))

	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	for _, element := range queries {
		tx.Exec(
			context.Background(),
			element.Sql,
			element.Params...,
		)
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func ParseInt(val string) *int {
	var result *int

	if val != "" {
		convertedVal, err := strconv.Atoi(val)
		if err != nil {
			result = &convertedVal
		}
	}

	return result
}

func ParseFloat64(val string) *float64 {
	var result *float64

	if val != "" {
		convertedVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			result = &convertedVal
		}
	}

	return result
}
