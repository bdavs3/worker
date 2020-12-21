package auth

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

// TODO (next): Only allow users to access endpoints for jobs they created.

const (
	storedUsername = "default_user"
	storedHash     = "$2a$10$P7GoVlD0fEu14OWE76dGzude2NLw0pi05Gzar6rm1b.oD04lcvyaq"
)

// TODO (out of scope): The userLimiter map should be cleaned up
// periodically to free memory. This could be done by keeping
// track of users' "last seen" times and having a background goroutine
// delete the oldest entries.
var userLimiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

// Secure enforces user authentication and rate limiting before allowing a
// request to reach a given endpoint.
func Secure(handler http.HandlerFunc) http.HandlerFunc {
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

		if isLocal, err := isLocalRequest(r); err != nil || !isLocal {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden non-local request."))
			return
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

func isLocalRequest(r *http.Request) (bool, error) {
	// TODO (out of scope): Getting the network address of requests using
	// RemoteAddr will unreliable they are first passed through a proxy or
	// load balancer. This should be taken into considereation.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false, err
	}

	reqIP := net.ParseIP(host)

	ifaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return false, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if reqIP.Equal(ip) {
				return true, nil
			}
		}
	}

	return false, nil
}
