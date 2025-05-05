package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthConfig holds the configuration for JWT authentication
type AuthConfig struct {
	SecretKey         string
	TokenExpiration   time.Duration
	RefreshExpiration time.Duration
}

// Claims represents the JWT token claims
type Claims struct {
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
	TokenID string `json:"token_id"` // Unique identifier for the token
	Type    string `json:"type"`     // "access" or "refresh"
	jwt.RegisteredClaims
}

// AuthUtils provides JWT authentication utilities
type AuthUtils struct {
	config AuthConfig
}

// NewAuthUtils creates a new AuthUtils instance
func NewAuthUtils(config AuthConfig) *AuthUtils {
	// Set default expiration times if not provided
	if config.TokenExpiration == 0 {
		config.TokenExpiration = 24 * time.Hour
	}
	if config.RefreshExpiration == 0 {
		config.RefreshExpiration = 7 * 24 * time.Hour
	}

	return &AuthUtils{
		config: config,
	}
}

// GenerateTokenID creates a unique token identifier
func GenerateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAccessToken creates a new JWT access token
func (au *AuthUtils) GenerateAccessToken(userID int, email string) (string, string, time.Time, error) {
	tokenID, err := GenerateTokenID()
	if err != nil {
		return "", "", time.Time{}, err
	}

	expiresAt := time.Now().Add(au.config.TokenExpiration)

	claims := &Claims{
		UserID:  userID,
		Email:   email,
		TokenID: tokenID,
		Type:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(au.config.SecretKey))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return signedToken, tokenID, expiresAt, nil
}

// GenerateRefreshToken creates a new JWT refresh token
func (au *AuthUtils) GenerateRefreshToken(userID int, email string) (string, string, time.Time, error) {
	tokenID, err := GenerateTokenID()
	if err != nil {
		return "", "", time.Time{}, err
	}

	expiresAt := time.Now().Add(au.config.RefreshExpiration)

	claims := &Claims{
		UserID:  userID,
		Email:   email,
		TokenID: tokenID,
		Type:    "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(au.config.SecretKey))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return signedToken, tokenID, expiresAt, nil
}

// VerifyToken validates a JWT token and returns the claims
func (au *AuthUtils) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(au.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// IsTokenExpired checks if a token has expired
func (au *AuthUtils) IsTokenExpired(claims *Claims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}

// ExtractTokenFromHeader gets the token from Authorization header
func ExtractTokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	const bearerPrefix = "Bearer "
	if len(header) < len(bearerPrefix) || header[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization header format")
	}

	return header[len(bearerPrefix):], nil
}

// TODO: Implement a proper password hashing and comparison mechanism

// HashPassword creates a password hash (you should use this for storing passwords)
func HashPassword(password string) (string, error) {
	// Use a proper password hashing library like bcrypt
	return password, nil // Replace with actual hashing implementation
}

// TODO: Implement a proper password comparison mechanism

// ComparePasswords compares a password with its hash
func ComparePasswords(password, hash string) error {
	// Use a proper password hashing library like bcrypt
	// This is just a placeholder
	if password == hash {
		return nil
	}
	return errors.New("password mismatch")
}
