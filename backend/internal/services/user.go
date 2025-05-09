package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bercivarga/website-builder/internal/models"
)

// UserService is a struct that holds the user store
type UserService struct {
	store *models.UserStore
}

// NewUserService creates a new UserService with the given UserStore
func NewUserService(store *models.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

// CreateUser handles the user registration process
func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if username == "" || email == "" || password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Username: username,
		Email:    email,
	}

	err := s.store.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// GetMe handles the retrieval of the current user
func (s *UserService) GetMe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetMe called")
	userID := r.Context().Value("userID").(int)
	fmt.Println("User ID from context:", userID)

	user, err := s.store.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Set Content-Type and encode JSON in one step
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		Username: user.Username,
		Email:    user.Email,
	})
}

// GetUser handles the retrieval of a user by ID
func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := s.store.GetUserByID(idInt)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	// write user data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
