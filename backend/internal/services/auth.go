package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bercivarga/website-builder/internal/models"
	"github.com/bercivarga/website-builder/internal/utils"
)

// AuthService is a struct that holds the token store and user store
type AuthService struct {
	store     *models.TokenStore
	authUtils *utils.AuthUtils
	userStore *models.UserStore
}

// UserCredentials represents login credentials
type UserCredentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenResponse represents the response after successful authentication
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// NewAuthService creates a new AuthService with the given stores and auth utils
func NewAuthService(store *models.TokenStore, authUtils *utils.AuthUtils, userStore *models.UserStore) *AuthService {
	return &AuthService{
		store:     store,
		authUtils: authUtils,
		userStore: userStore,
	}
}

// AuthMiddleware is a middleware for protecting routes
func (s *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		tokenString, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify token
		claims, err := s.authUtils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if token type is access
		if claims.Type != "access" {
			http.Error(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		// Check if token exists in database
		tokenDB, err := s.store.GetTokenByUserID(claims.UserID)
		if err != nil {
			http.Error(w, "Token not found", http.StatusUnauthorized)
			return
		}

		// Verify token ID matches
		if tokenDB.Token != claims.TokenID {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims and userId to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Login handles the user login process
func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var creds UserCredentials

	// Parse JSON body
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate user credentials
	user, err := s.userStore.GetUserByEmail(creds.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Verify password
	err = utils.ComparePasswords(creds.Password, user.PasswordHash)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate access token
	accessToken, accessTokenID, accessExpires, err := s.authUtils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken, refreshTokenID, refreshExpires, err := s.authUtils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Store access token in database
	accessTokenDB := &models.Token{
		UserID:    user.ID,
		Token:     accessTokenID,
		TokenType: "access",
		ExpiresAt: accessExpires,
	}

	err = s.store.CreateToken(accessTokenDB)
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		return
	}

	// Store refresh token in database
	refreshTokenDB := &models.Token{
		UserID:    user.ID,
		Token:     refreshTokenID,
		TokenType: "refresh",
		ExpiresAt: refreshExpires,
	}

	err = s.store.CreateToken(refreshTokenDB)
	if err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	// Send response
	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpires.Sub(time.Now()).Seconds()),
		TokenType:    "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register handles user registration
func (s *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	var creds UserCredentials

	// Parse JSON body
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate user credentials
	if creds.Email == "" || creds.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Hash password
	passwordHash, err := utils.HashPassword(creds.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user in database
	user := &models.User{
		Email:        creds.Email,
		PasswordHash: passwordHash,
		Username:     creds.Username,
	}

	err = s.userStore.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// Refresh handles token refresh requests
func (s *AuthService) Refresh(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from Authorization header
	authHeader := r.Header.Get("Authorization")
	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	// Verify refresh token
	claims, err := s.authUtils.VerifyToken(tokenString)
	if err != nil || claims.Type != "refresh" {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Check if token exists in database
	tokenDB, err := s.store.GetTokenByUserID(claims.UserID)
	if err != nil {
		http.Error(w, "Token not found", http.StatusUnauthorized)
		return
	}

	// Verify token ID matches
	if tokenDB.Token != claims.TokenID {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Generate new access token
	accessToken, accessTokenID, accessExpires, err := s.authUtils.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Store new access token
	accessTokenDB := &models.Token{
		UserID:    claims.UserID,
		Token:     accessTokenID,
		TokenType: "access",
		ExpiresAt: accessExpires,
	}

	err = s.store.CreateToken(accessTokenDB)
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		return
	}

	// Send response
	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: tokenString, // Keep the same refresh token
		ExpiresIn:    int64(accessExpires.Sub(time.Now()).Seconds()),
		TokenType:    "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout by revoking tokens
func (s *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	// Extract access token from Authorization header
	authHeader := r.Header.Get("Authorization")
	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	// Verify access token
	claims, err := s.authUtils.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Revoke all tokens for this user
	err = s.store.DeleteTokensByUserID(claims.UserID)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
