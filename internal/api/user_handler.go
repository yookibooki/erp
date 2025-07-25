package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yookibooki/erp/internal/auth"
	"github.com/yookibooki/erp/internal/models"
)

// UserHandler handles user requests
type UserHandler struct {
	userService models.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService models.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser gets a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	user, err := h.userService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting user")
		return
	}

	if user == nil {
		auth.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, user)
}

// ListUsers lists all users for a tenant
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	users, err := h.userService.List(tenantID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing users")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, users)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set tenant ID from context
	user.TenantID = tenantID

	// Validate user
	if user.Email == "" || user.PasswordHash == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Check if user already exists
	existingUser, err := h.userService.GetByEmail(tenantID, user.Email)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking user")
		return
	}

	if existingUser != nil {
		auth.RespondWithError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	// Create user
	if err := h.userService.Create(&user); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, user)
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and tenant ID
	user.ID = id
	user.TenantID = tenantID

	// Validate user
	if user.Email == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Check if user exists
	existingUser, err := h.userService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking user")
		return
	}

	if existingUser == nil {
		auth.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Update user
	if err := h.userService.Update(&user); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	// Check if user exists
	existingUser, err := h.userService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking user")
		return
	}

	if existingUser == nil {
		auth.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Delete user
	if err := h.userService.Delete(tenantID, id); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error deleting user")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}