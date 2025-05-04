package handlers

import (
	"net/http"

	"github.com/bercivarga/website-builder/internal/app"
)

// SetupHandlers initializes the HTTP routes for the application.
func SetupHandlers(app *app.Application) *http.ServeMux {
	mux := http.NewServeMux()

	addPublicRoutes(mux)
	addAuthRoutes(mux, app)
	addUserRoutes(mux, app)

	return mux
}

func addPublicRoutes(mux *http.ServeMux) {
	publicGroup := CreateRouteGroup(mux, "/v1")
	publicGroup.Use(LoggingMiddleware)

	publicGroup.Get("/health", healthCheckHandler)
}

func addAuthRoutes(mux *http.ServeMux, app *app.Application) {
	authGroup := CreateRouteGroup(mux, "/v1/auth")
	authGroup.Use(LoggingMiddleware)

	authGroup.Post("/register", registerUserHandler)
	authGroup.Post("/login", loginUserHandler)
}

func addUserRoutes(mux *http.ServeMux, app *app.Application) {
	userGroup := CreateRouteGroup(mux, "/v1/user")
	userGroup.Use(LoggingMiddleware)
	userGroup.Use(AuthMiddleware)

	userGroup.Get("/{id}", app.UserService.GetUser)
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User registration"))
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User login"))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
