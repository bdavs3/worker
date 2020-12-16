package auth

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser verifies the credentials in the Authorization header of a request by utilizing the bcrypt hashing function.
func AuthenticateUser(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, pw, ok := r.BasicAuth(); !ok || !validate(username, pw) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Account Invalid"`)
			w.WriteHeader(401)
			w.Write([]byte("Invalid credentials. Access denied."))
		} else {
			handler(w, r)
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
