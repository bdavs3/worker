package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	storedUsername = "default_user"
	storedHash     = "$2a$10$P7GoVlD0fEu14OWE76dGzude2NLw0pi05Gzar6rm1b.oD04lcvyaq"
)

// A SecurityLayer can perform security checks on HTTP handlers and set owner
// relationships to id-represented resources.
type SecurityLayer interface {
	Authenticate(handler http.Handler) http.Handler
	Authorize(handler http.Handler) http.Handler
	SetOwner(username, id string)
}

// DummyAuth is a SecurityLayer intended only for testing dependent functions.
type DummyAuth struct{}

func (da *DummyAuth) Authenticate(handler http.Handler) http.Handler { return nil }
func (da *DummyAuth) Authorize(handler http.Handler) http.Handler    { return nil }
func (da *DummyAuth) SetOwner(username, id string)                   {}

// Auth is used to enforce security checks on client requests. Use NewAuth to create a
// new instance.
type Auth struct {
	owners *ownershipTracker
}

// NewAuth creates a new instance of the auth layer.
func NewAuth() *Auth {
	return &Auth{
		owners: newOwnershipTracker(),
	}
}

// Authenticate performs an authentication check on an HTTP Handler.
func (a *Auth) Authenticate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, pw, ok := r.BasicAuth()

		if !ok || !validate(username, pw) {
			http.Error(w, "invalid credentials: access denied", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
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

// Authorize performs a resource-ownership check on an HTTP handler.
func (a *Auth) Authorize(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, ok := r.BasicAuth()

		id := mux.Vars(r)["id"]
		if !ok || !a.owners.isOwner(username, id) {
			// If a user tries to access an endpoint belonging to someone else, do not
			// reveal that the endpoint exists by responding with StatusNotFound.
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// SetOwner designates the given user as the owner of the resource with the given id.
func (a *Auth) SetOwner(username, id string) {
	// TODO (out of scope): Track owner relationships in a database.
	a.owners.setOwner(username, id)
}
