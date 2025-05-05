package handlers

import (
	"net/http"

	"github.com/bercivarga/website-builder/internal/app"
	"github.com/rs/cors"
)

// SetupHandlers initializes the HTTP routes for the application.
func SetupHandlers(app *app.Application) http.Handler {
	mux := http.NewServeMux()
	addPublicRoutes(mux, app)
	addAuthRoutes(mux, app)
	addUserRoutes(mux, app)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "your-frontend-url"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
		// Debug:            true, // Enable for debugging
	})

	return c.Handler(mux)
}

func addPublicRoutes(mux *http.ServeMux, app *app.Application) {
	publicGroup := CreateRouteGroup(mux, "/v1")
	publicGroup.Use(LoggingMiddleware(app.Logger))
	publicGroup.Get("/health", healthCheckHandler)
}

func addAuthRoutes(mux *http.ServeMux, app *app.Application) {
	authGroup := CreateRouteGroup(mux, "/v1/auth")
	authGroup.Use(LoggingMiddleware(app.Logger))
	authGroup.Post("/register", app.AuthService.Register)
	authGroup.Post("/login", app.AuthService.Login)
	authGroup.Post("/logout", app.AuthService.Logout)
	authGroup.Post("/refresh", app.AuthService.Refresh)
}

func addUserRoutes(mux *http.ServeMux, app *app.Application) {
	userGroup := CreateRouteGroup(mux, "/v1/user")
	userGroup.Use(LoggingMiddleware(app.Logger))
	userGroup.Use(app.AuthService.AuthMiddleware)
	userGroup.Get("/me", app.UserService.GetMe)
	userGroup.Get("/{id}", app.UserService.GetUser)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
