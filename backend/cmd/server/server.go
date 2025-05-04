package server

import (
	"net/http"
	"time"

	"github.com/bercivarga/website-builder/internal/app"
	"github.com/bercivarga/website-builder/internal/handlers"
)

const (
	port = ":8080"
)

// Start initializes the application and starts the HTTP server.
func Start() (*app.Application, error) {
	app, err := app.NewApplication()
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:         port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handlers.SetupHandlers(&app),
	}

	err = server.ListenAndServe()

	if err != nil {
		return nil, err
	}

	return &app, nil
}
