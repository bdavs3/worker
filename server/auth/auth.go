package auth

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(5, 1) // Allows a request every 200ms.

// Secure enforces user authentication and rate limiting before allowing a request to reach a given endpoint.
func Secure(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, pw, ok := r.BasicAuth(); !ok || !validate(username, pw) {
			w.WriteHeader(401)
			w.Write([]byte("Invalid credentials. Access denied."))
		} else if !limiter.Allow() { // TODO (next): Enforce per-user.
			w.WriteHeader(429)
			w.Write([]byte("Too many requests."))
		} else {
			router.ServeHTTP(w, r)
		}
	}
}

func validate(username, pw string) bool {
	// TODO (out of scope): Store user credentials in a secure database and
	// validate request Authorization headers against them. It is critical
	// that passwords are hashed before storage in the database.
	if username == "default_user" {
		// bcrypt cost of 10 is chosen because it takes roughly 75-100ms to
		// execute, a benchmark that is relatively unnoticeable to the user
		// but resiliant against brute-force attacks.
		hash, err := bcrypt.GenerateFromPassword([]byte("123456"), 10)
		if err != nil {
			log.Fatal(err)
		}

		err = bcrypt.CompareHashAndPassword(hash, []byte(pw))
		return err == nil
	}
	return false
}
