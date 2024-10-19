package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

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
