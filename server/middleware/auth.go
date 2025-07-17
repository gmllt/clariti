package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gmllt/clariti/logger"
	"github.com/gmllt/clariti/server/config"
)

// BasicAuth provides HTTP Basic Authentication middleware
func BasicAuth(config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetDefault().WithComponent("BasicAuthMiddleware")

			// Skip auth for GET requests (read-only)
			if r.Method == http.MethodGet {
				log.WithField("method", r.Method).WithField("path", r.URL.Path).Debug("Skipping auth for read-only request")
				next.ServeHTTP(w, r)
				return
			}

			log.WithField("method", r.Method).WithField("path", r.URL.Path).Debug("Checking basic auth for write request")

			username, password, ok := r.BasicAuth()
			if !ok {
				log.Warn("Basic auth credentials missing")
				w.Header().Set("WWW-Authenticate", `Basic realm="Clariti API"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Use constant-time comparison to prevent timing attacks
			usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(config.Auth.AdminUsername)) == 1
			passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(config.Auth.AdminPassword)) == 1

			if !usernameMatch || !passwordMatch {
				log.WithField("username", username).Warn("Authentication failed - invalid credentials")
				w.Header().Set("WWW-Authenticate", `Basic realm="Clariti API"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			log.WithField("username", username).Info("Authentication successful")
			next.ServeHTTP(w, r)
		})
	}
}

// CORS adds basic CORS headers
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetDefault().WithComponent("CORSMiddleware")

		log.WithField("method", r.Method).WithField("origin", r.Header.Get("Origin")).Debug("Processing CORS headers")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			log.WithField("path", r.URL.Path).Debug("Handling preflight CORS request")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
