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

// A UserAuthLayer authenticates/authorizes users by validating their credentials
// and verifying that they are owner of any resource they are attempting to
// access. It also enables users to be set as owner of a particular resource.
type UserAuthLayer interface {
	Secure(handler http.Handler) http.Handler
	SetOwner(username, id string)
}

// DummyAuth implements the AuthLayer interface so that the API can be tested
// independently.
type DummyAuth struct{}

func (da *DummyAuth) Secure(handler http.Handler) http.Handler { return nil }
func (da *DummyAuth) SetOwner(username, id string)             {}

// Auth provides the ability to authenticate/authorize users and to set them
// as owners of particular resources. Use NewAuth to create a new instance.
type Auth struct {
	owners *ownershipTracker
}

// NewAuth returns an Auth layer with an empty ownership map.
func NewAuth() *Auth {
	return &Auth{
		owners: newOwnershipTracker(),
	}
}

// Secure enforces user authentication and a resource ownership check on
// an HTTP Handler.
func (a *Auth) Secure(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, pw, ok := r.BasicAuth()

		if !ok || !validate(username, pw) {
			http.Error(w, "invalid credentials: access denied", http.StatusUnauthorized)
			return
		}

		if r.Method != http.MethodPost {
			id := mux.Vars(r)["id"]
			if !a.owners.isOwner(username, id) {
				// If a user tries to access an endpoint belonging to someone else, do not
				// reveal that the endpoint exists by responding with StatusNotFound.
				http.Error(w, "job not found", http.StatusNotFound)
				return
			}
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

// SetOwner designates the given user as the owner of the resource with the
// given id.
func (a *Auth) SetOwner(username, id string) {
	a.owners.setOwner(username, id)
}
