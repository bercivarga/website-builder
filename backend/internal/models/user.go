package models

import "database/sql"

// User represents a user in the system
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// UserStore is a struct that holds the database connection
type UserStore struct {
	DB *sql.DB
}

// UserRepository is an interface that defines the methods for user operations
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
}

// NewUserStore creates a new UserStore with the given database connection
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{DB: db}
}

// CreateUser inserts a new user into the database
func (s *UserStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.DB.Exec(query, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by ID from the database
func (s *UserStore) GetUserByID(id int) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	row := s.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username from the database
func (s *UserStore) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1`
	row := s.DB.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database
func (s *UserStore) UpdateUser(user *User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, updated_at = $4 WHERE id = $5`
	_, err := s.DB.Exec(query, user.Username, user.Email, user.PasswordHash, user.UpdatedAt, user.ID)
	return err
}

// DeleteUser deletes a user from the database
func (s *UserStore) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}
