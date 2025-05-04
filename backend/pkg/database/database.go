package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Connect establishes a connection to the PostgreSQL database using environment variables.
func Connect() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return nil, fmt.Errorf("missing one or more required environment variables")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60) // 5 minutes
	db.SetConnMaxIdleTime(5 * 60) // 5 minutes

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Log the successful connection
	fmt.Println("Connected to the database successfully!")

	return db, nil
}
