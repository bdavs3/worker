package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// TODO (next): Only allow users to access endpoints for jobs they created.

const (
	storedUsername = "default_user"
	storedHash     = "$2a$10$P7GoVlD0fEu14OWE76dGzude2NLw0pi05Gzar6rm1b.oD04lcvyaq"
)

type Auth struct {
	owners *ownershipTracker
}

func NewAuth() *Auth {
	return &Auth{
		owners: newOwnershipTracker(),
	}
}

// Secure enforces user authentication.
func (a *Auth) Secure(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, pw, ok := r.BasicAuth()

		if !ok || !validate(username, pw) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid credentials. Access denied."))
			return
		}

		if r.Method != http.MethodPost {
			id := mux.Vars(r)["id"]
			if !a.owners.
		}

		handler(w, r)
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
