package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// BasicAuth HTTP Basic Authentication middleware
type BasicAuth struct {
	username string
	password string
}

// NewBasicAuth creates a new instance of BasicAuth middleware
func NewBasicAuth() *BasicAuth {
	return &BasicAuth{
		username: os.Getenv("BASIC_AUTH_USERNAME"),
		password: os.Getenv("BASIC_AUTH_PASSWORD"),
	}
}

// IsEnabled checks if basic authentication is enabled
func (ba *BasicAuth) IsEnabled() bool {
	return ba.username != "" && ba.password != ""
}

// Middleware wraps http.HandlerFunc with basic authentication
func (ba *BasicAuth) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ba.IsEnabled() {
			next(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			ba.unauthorized(w)
			return
		}

		usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(ba.username)) == 1
		passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(ba.password)) == 1

		if !usernameMatch || !passwordMatch {
			ba.unauthorized(w)
			return
		}

		next(w, r)
	}
}

func (ba *BasicAuth) unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"Unauthorized"}`))
}
