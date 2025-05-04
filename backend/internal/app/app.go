package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/bercivarga/website-builder/internal/models"
	"github.com/bercivarga/website-builder/internal/services"
	"github.com/bercivarga/website-builder/migrations"
	"github.com/bercivarga/website-builder/pkg/database"
	"github.com/joho/godotenv"
)

// Application holds the application state
type Application struct {
	DB          *sql.DB
	Logger      *log.Logger
	UserService *services.UserService
}

// NewApplication initializes the application with a database connection and logger.
func NewApplication() (Application, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return Application{}, err
	}

	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		return Application{}, err
	}

	// Migrate the database
	err = database.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return Application{}, err
	}

	// Initialize logger
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	if logger == nil {
		return Application{}, os.ErrInvalid
	}

	// stores go here
	userStore := models.NewUserStore(db)

	// services go here
	userService := services.NewUserService(userStore)

	app := Application{
		DB:          db,
		Logger:      logger,
		UserService: userService,
	}

	return app, nil
}
