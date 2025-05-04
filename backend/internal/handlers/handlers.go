package handlers

import (
	"net/http"

	"github.com/bercivarga/website-builder/internal/app"
)

// SetupHandlers initializes the HTTP routes for the application.
func SetupHandlers(app *app.Application) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check route
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// User routes
	userGroup := CreateRouteGroup(mux, "/v1/user")
	userGroup.Use(LoggingMiddleware)
	userGroup.Use(AuthMiddleware)

	userGroup.Post("/register", registerUserHandler)
	userGroup.Post("/login", loginUserHandler)

	return mux
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user registration
	w.Write([]byte("User registration"))
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user login
	w.Write([]byte("User login"))
}
