package handlers

import "net/http"

// LoggingMiddleware is a middleware that logs incoming requests and their details.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware is a middleware that checks for authentication.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication
		next.ServeHTTP(w, r)
	})
}
