package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yookibooki/erp/internal/auth"
	"github.com/yookibooki/erp/internal/models"
)

// TenantHandler handles tenant requests
type TenantHandler struct {
	tenantService models.TenantService
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantService models.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// GetTenant gets a tenant by ID
func (h *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tenant, err := h.tenantService.GetByID(id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting tenant")
		return
	}

	if tenant == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Tenant not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, tenant)
}

// GetTenantBySubdomain gets a tenant by subdomain
func (h *TenantHandler) GetTenantBySubdomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subdomain := vars["subdomain"]

	tenant, err := h.tenantService.GetBySubdomain(subdomain)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting tenant")
		return
	}

	if tenant == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Tenant not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, tenant)
}

// ListTenants lists all tenants
func (h *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.tenantService.List()
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing tenants")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, tenants)
}

// CreateTenant creates a new tenant
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var tenant models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate tenant
	if tenant.Name == "" || tenant.Subdomain == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Name and subdomain are required")
		return
	}

	// Check if tenant already exists
	existingTenant, err := h.tenantService.GetBySubdomain(tenant.Subdomain)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking tenant")
		return
	}

	if existingTenant != nil {
		auth.RespondWithError(w, http.StatusConflict, "Tenant with this subdomain already exists")
		return
	}

	// Create tenant
	if err := h.tenantService.Create(&tenant); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating tenant")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, tenant)
}

// UpdateTenant updates a tenant
func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var tenant models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID from URL
	tenant.ID = id

	// Validate tenant
	if tenant.Name == "" || tenant.Subdomain == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Name and subdomain are required")
		return
	}

	// Check if tenant exists
	existingTenant, err := h.tenantService.GetByID(id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking tenant")
		return
	}

	if existingTenant == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Tenant not found")
		return
	}

	// Update tenant
	if err := h.tenantService.Update(&tenant); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error updating tenant")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, tenant)
}

// DeleteTenant deletes a tenant
func (h *TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if tenant exists
	existingTenant, err := h.tenantService.GetByID(id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking tenant")
		return
	}

	if existingTenant == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Tenant not found")
		return
	}

	// Delete tenant
	if err := h.tenantService.Delete(id); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error deleting tenant")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Tenant deleted successfully"})
}