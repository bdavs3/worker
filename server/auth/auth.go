package auth

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// TODO (next): Only allow users to access endpoints for jobs they created.

const (
	storedUsername = "default_user"
	storedHash     = "$2a$10$P7GoVlD0fEu14OWE76dGzude2NLw0pi05Gzar6rm1b.oD04lcvyaq"
)

// Secure enforces user authentication on an HTTP Handler.
func Secure(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, pw, ok := r.BasicAuth()

		if !ok || !validate(username, pw) {
			http.Error(w, "invalid credentials: access denied", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	}
}

func validate(username, pw string) bool {
	// TODO (out of scope): Store user credentials in a secure database and
	// validate request Authorization headers against them. It is critical
	// that passwords are hashed before storage in the database.
	if username == storedUsername {
		err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(pw))
		return err == nil
	}
	return false
}
