package middleware

import (
	"net/http"
	"time"

	"github.com/gmllt/clariti/logger"
	"github.com/gmllt/clariti/server/metrics"
	"github.com/sirupsen/logrus"
)

// ResponseWriter wrapper to capture status code - flight recorder instrument
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

// MetricsMiddleware captures HTTP metrics for monitoring - flight data collection
func MetricsMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				written:        false,
			}

			// Process request
			next.ServeHTTP(wrapped, r)

			// Record metrics
			duration := time.Since(start)
			metrics.RecordHTTPRequest(r.Method, r.URL.Path, wrapped.statusCode, duration)

			// Log with aviation terminology
			log.WithFields(logrus.Fields{
				"component":   "MetricsMiddleware",
				"method":      r.Method,
				"path":        r.URL.Path,
				"status_code": wrapped.statusCode,
				"duration_ms": duration.Milliseconds(),
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			}).Debug("Flight data recorded - HTTP request metrics captured")
		})
	}
}
