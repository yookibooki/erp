package api

import (
	"encoding/json"
	"net/http"

	"github.com/yookibooki/erp/internal/auth"
	"github.com/yookibooki/erp/internal/models"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService models.UserService
	jwtService  *auth.JWTService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService models.UserService, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Login handles login requests
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate request
	if req.TenantID == "" || req.Email == "" || req.Password == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Tenant ID, email and password are required")
		return
	}

	// Authenticate user
	user, err := h.userService.Authenticate(req.TenantID, req.Email, req.Password)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error authenticating user")
		return
	}

	if user == nil {
		auth.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Return token and user
	auth.RespondWithJSON(w, http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// RegisterRequest represents a register request
type RegisterRequest struct {
	TenantID  string `json:"tenant_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// Register handles register requests
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate request
	if req.TenantID == "" || req.Email == "" || req.Password == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Tenant ID, email and password are required")
		return
	}

	// Check if user already exists
	existingUser, err := h.userService.GetByEmail(req.TenantID, req.Email)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking user")
		return
	}

	if existingUser != nil {
		auth.RespondWithError(w, http.StatusConflict, "User already exists")
		return
	}

	// Create user
	user := &models.User{
		TenantID:     req.TenantID,
		Email:        req.Email,
		PasswordHash: req.Password, // Will be hashed in the repository
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
	}

	if err := h.userService.Create(user); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Return token and user
	auth.RespondWithJSON(w, http.StatusCreated, LoginResponse{
		Token: token,
		User:  user,
	})
}