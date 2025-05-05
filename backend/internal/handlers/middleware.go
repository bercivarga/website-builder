package handlers

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware is a middleware that logs incoming requests and their details.
func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the ResponseWriter to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			// Log request details
			duration := time.Since(start)
			logger.Printf(
				"%s %s %s %d %v",
				r.RemoteAddr,
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
			)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
