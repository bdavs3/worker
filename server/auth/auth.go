package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

// TODO (next): Enhance the rate-limiting authorization check by making it
// user-specific. Add two additional authorization checks: (1) Only allow
// users to access endpoints for jobs they created and (2) Enforce local-only
// request origin policy.

const (
	globalUsername = "default_user"
	globalHash     = "$2a$10$P7GoVlD0fEu14OWE76dGzude2NLw0pi05Gzar6rm1b.oD04lcvyaq"
)

var limiter = rate.NewLimiter(5, 1) // Allows a request every 200ms.

// Secure enforces user authentication and rate limiting before allowing a
// request to reach a given endpoint.
func Secure(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, pw, ok := r.BasicAuth(); !ok || !validate(username, pw) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid credentials. Access denied."))
			return
		}
		if !limiter.Allow() { // TODO (next): Enforce per-user.
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests."))
			return
		}
		router.ServeHTTP(w, r)
	}
}

func validate(username, pw string) bool {
	// TODO (out of scope): Store user credentials in a secure database and
	// validate request Authorization headers against them. It is critical
	// that passwords are hashed before storage in the database.
	if username == globalUsername {
		err := bcrypt.CompareHashAndPassword([]byte(globalHash), []byte(pw))
		return err == nil
	}
	return false
}
