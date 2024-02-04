package pkg

import (
	"net/http"
	"strings"
)

var jwtSecret = []byte("Asdf@123")
var hardcodedToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For demonstration purposes, used the hardcoded token
		tokenString := hardcodedToken

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Unauthorized - Missing token", http.StatusUnauthorized)
			return
		}

		if !strings.EqualFold(tokenHeader, "Bearer "+tokenString) {
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
