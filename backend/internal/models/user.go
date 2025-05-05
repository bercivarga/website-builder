package models

import (
	"database/sql"
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't include in JSON
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserStore is a struct that holds the database connection
type UserStore struct {
	DB *sql.DB
}

// UserRepository is an interface that defines the methods for user operations
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
}

// NewUserStore creates a new UserStore with the given database connection
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{DB: db}
}

// CreateUser inserts a new user into the database
func (s *UserStore) CreateUser(user *User) error {
	query := `INSERT INTO users (email, password_hash, username) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := s.DB.QueryRow(query, user.Email, user.PasswordHash, user.Username).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

// GetUserByID retrieves a user by ID from the database
func (s *UserStore) GetUserByID(id int) (*User, error) {
	query := `SELECT id, email, password_hash, username, created_at, updated_at FROM users WHERE id = $1`
	row := s.DB.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email from the database
func (s *UserStore) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password_hash, username, created_at, updated_at FROM users WHERE email = $1`
	row := s.DB.QueryRow(query, email)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database
func (s *UserStore) UpdateUser(user *User) error {
	query := `UPDATE users SET email = $1, password_hash = $2, username = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := s.DB.Exec(query, user.Email, user.PasswordHash, user.Username, user.ID)
	return err
}

// DeleteUser deletes a user from the database
func (s *UserStore) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}
