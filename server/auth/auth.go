package auth

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

// TODO (next): Add two additional authorization checks: (1) Only allow
// users to access endpoints for jobs they created and (2) Enforce local-only
// request origin policy.

const (
	globalUsername = "default_user"
	globalHash     = "$2a$10$P7GoVlD0fEu14OWE76dGzude2NLw0pi05Gzar6rm1b.oD04lcvyaq"
)

// TODO (out of scope): The userLimiter map should be cleaned up
// periodically to free memory. This could be done by keeping
// track of users' "last seen" times and having a background goroutine
// delete the oldest entries.
var userLimiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

// Secure enforces user authentication and rate limiting before allowing a
// request to reach a given endpoint.
func Secure(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, pw, ok := r.BasicAuth()

		if !ok || !validate(username, pw) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid credentials. Access denied."))
			return
		}

		if !getUserLimiter(username).Allow() {
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

func getUserLimiter(username string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := userLimiters[username]
	if !exists {
		limiter = rate.NewLimiter(5, 1) // Allows a request every 200ms.
		userLimiters[username] = limiter
	}

	return limiter
}
