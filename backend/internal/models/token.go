package models

import (
	"database/sql"
	"time"
)

// Token represents a user authentication token
type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"` // "access" or "refresh"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TokenStore is a struct that holds the database connection
type TokenStore struct {
	DB *sql.DB
}

// TokenRepository is an interface that defines the methods for token operations
type TokenRepository interface {
	CreateToken(token *Token) error
	GetTokenByID(id int) (*Token, error)
	GetTokenByUserID(userID int) (*Token, error)
	GetTokensByUserID(userID int) ([]*Token, error)
	DeleteToken(id int) error
	DeleteTokensByUserID(userID int) error
	UpdateToken(token *Token) error
}

// NewTokenStore creates a new TokenStore with the given database connection
func NewTokenStore(db *sql.DB) *TokenStore {
	return &TokenStore{DB: db}
}

// CreateToken inserts a new token into the database
func (s *TokenStore) CreateToken(token *Token) error {
	query := `INSERT INTO token (user_id, token, token_type, expires_at) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := s.DB.QueryRow(query, token.UserID, token.Token, token.TokenType, token.ExpiresAt).Scan(&token.ID, &token.CreatedAt, &token.UpdatedAt)
	return err
}

// GetTokenByID retrieves a token by ID from the database
func (s *TokenStore) GetTokenByID(id int) (*Token, error) {
	query := `SELECT id, user_id, token, token_type, created_at, updated_at, expires_at FROM token WHERE id = $1`
	row := s.DB.QueryRow(query, id)
	var token Token
	err := row.Scan(&token.ID, &token.UserID, &token.Token, &token.TokenType, &token.CreatedAt, &token.UpdatedAt, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetTokenByUserID retrieves a token by user ID from the database
func (s *TokenStore) GetTokenByUserID(userID int) (*Token, error) {
	query := `SELECT id, user_id, token, token_type, created_at, updated_at, expires_at FROM token WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	row := s.DB.QueryRow(query, userID)
	var token Token
	err := row.Scan(&token.ID, &token.UserID, &token.Token, &token.TokenType, &token.CreatedAt, &token.UpdatedAt, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetTokensByUserID retrieves all tokens for a user from the database
func (s *TokenStore) GetTokensByUserID(userID int) ([]*Token, error) {
	query := `SELECT id, user_id, token, token_type, created_at, updated_at, expires_at FROM token WHERE user_id = $1`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*Token
	for rows.Next() {
		var token Token
		err := rows.Scan(&token.ID, &token.UserID, &token.Token, &token.TokenType, &token.CreatedAt, &token.UpdatedAt, &token.ExpiresAt)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}
	return tokens, nil
}

// DeleteToken removes a token from the database
func (s *TokenStore) DeleteToken(id int) error {
	query := `DELETE FROM token WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

// DeleteTokensByUserID removes all tokens for a user from the database
func (s *TokenStore) DeleteTokensByUserID(userID int) error {
	query := `DELETE FROM token WHERE user_id = $1`
	_, err := s.DB.Exec(query, userID)
	return err
}

// UpdateToken updates an existing token in the database
func (s *TokenStore) UpdateToken(token *Token) error {
	query := `UPDATE token SET user_id = $1, token = $2, token_type = $3, expires_at = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	_, err := s.DB.Exec(query, token.UserID, token.Token, token.TokenType, token.ExpiresAt, token.ID)
	return err
}
