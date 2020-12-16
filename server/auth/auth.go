package auth

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(5, 1) // Allows request every 200ms.

// Secure enforces user authentication and rate limiting before allowing the router to direct a request to a given endpoint.
func Secure(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, pw, ok := r.BasicAuth(); !ok || !validate(username, pw) {
			w.WriteHeader(401)
			w.Write([]byte("Invalid credentials. Access denied."))
		} else if !limiter.Allow() { // TODO: Enforce per-user
			w.WriteHeader(429)
			w.Write([]byte("Too many requests."))
		} else {
			router.ServeHTTP(w, r)
		}
	}
}

func validate(username, pw string) bool {
	if username == "default_user" {
		// TODO: Store user credentials in a secure database.
		// User passwords should be hashed before storing them in the DB.
		hash, err := hashPassword("123456")
		if err != nil {
			log.Fatal(err)
		}

		return checkPasswordHash(pw, hash)
	}
	return false
}

func hashPassword(password string) (string, error) {
	// bcrypt cost of 10 is chosen because it takes roughly 75-100ms to execute, a benchmark that is relatively unnoticeable to the user but resiliant against brute-force attacks.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
