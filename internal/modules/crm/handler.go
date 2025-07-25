package crm

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yookibooki/erp/internal/auth"
	"github.com/yookibooki/erp/internal/models"
)

// CustomerHandler handles customer requests
type CustomerHandler struct {
	customerService models.CustomerService
	contactService  models.ContactService
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(
	customerService models.CustomerService,
	contactService models.ContactService,
) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		contactService:  contactService,
	}
}

// GetCustomer gets a customer by ID
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	customer, err := h.customerService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting customer")
		return
	}

	if customer == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	// Get contacts for customer
	contacts, err := h.contactService.ListByCustomer(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting contacts")
		return
	}

	// Convert []*models.Contact to []models.Contact
	customer.Contacts = make([]models.Contact, len(contacts))
	for i, contact := range contacts {
		customer.Contacts[i] = *contact
	}

	auth.RespondWithJSON(w, http.StatusOK, customer)
}

// ListCustomers lists all customers for a tenant
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	customers, err := h.customerService.List(tenantID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing customers")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, customers)
}

// CreateCustomer creates a new customer
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set tenant ID from context
	customer.TenantID = tenantID

	// Validate customer
	if customer.Name == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Name is required")
		return
	}

	// Create customer
	if err := h.customerService.Create(&customer); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating customer")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, customer)
}

// UpdateCustomer updates a customer
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and tenant ID
	customer.ID = id
	customer.TenantID = tenantID

	// Validate customer
	if customer.Name == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Name is required")
		return
	}

	// Check if customer exists
	existingCustomer, err := h.customerService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking customer")
		return
	}

	if existingCustomer == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	// Update customer
	if err := h.customerService.Update(&customer); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error updating customer")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, customer)
}

// DeleteCustomer deletes a customer
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	// Check if customer exists
	existingCustomer, err := h.customerService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking customer")
		return
	}

	if existingCustomer == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	// Delete customer
	if err := h.customerService.Delete(tenantID, id); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error deleting customer")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Customer deleted successfully"})
}

// ContactHandler handles contact requests
type ContactHandler struct {
	contactService models.ContactService
	customerService models.CustomerService
}

// NewContactHandler creates a new contact handler
func NewContactHandler(
	contactService models.ContactService,
	customerService models.CustomerService,
) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
		customerService: customerService,
	}
}

// GetContact gets a contact by ID
func (h *ContactHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	contact, err := h.contactService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting contact")
		return
	}

	if contact == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Contact not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, contact)
}

// ListContactsByCustomer lists all contacts for a customer
func (h *ContactHandler) ListContactsByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	contacts, err := h.contactService.ListByCustomer(tenantID, customerID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing contacts")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, contacts)
}

// CreateContact creates a new contact
func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set tenant ID from context
	contact.TenantID = tenantID

	// Validate contact
	if contact.CustomerID == "" || contact.FirstName == "" || contact.LastName == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Customer ID, first name and last name are required")
		return
	}

	// Check if customer exists
	customer, err := h.customerService.GetByID(tenantID, contact.CustomerID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking customer")
		return
	}

	if customer == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	// Create contact
	if err := h.contactService.Create(&contact); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating contact")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, contact)
}

// UpdateContact updates a contact
func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and tenant ID
	contact.ID = id
	contact.TenantID = tenantID

	// Validate contact
	if contact.CustomerID == "" || contact.FirstName == "" || contact.LastName == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Customer ID, first name and last name are required")
		return
	}

	// Check if contact exists
	existingContact, err := h.contactService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking contact")
		return
	}

	if existingContact == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Contact not found")
		return
	}

	// Update contact
	if err := h.contactService.Update(&contact); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error updating contact")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, contact)
}

// DeleteContact deletes a contact
func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	// Check if contact exists
	existingContact, err := h.contactService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking contact")
		return
	}

	if existingContact == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Contact not found")
		return
	}

	// Delete contact
	if err := h.contactService.Delete(tenantID, id); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error deleting contact")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Contact deleted successfully"})
}

// InteractionHandler handles interaction requests
type InteractionHandler struct {
	interactionService models.InteractionService
	customerService    models.CustomerService
}

// NewInteractionHandler creates a new interaction handler
func NewInteractionHandler(
	interactionService models.InteractionService,
	customerService models.CustomerService,
) *InteractionHandler {
	return &InteractionHandler{
		interactionService: interactionService,
		customerService:    customerService,
	}
}

// GetInteraction gets an interaction by ID
func (h *InteractionHandler) GetInteraction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	interaction, err := h.interactionService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting interaction")
		return
	}

	if interaction == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Interaction not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, interaction)
}

// ListInteractionsByCustomer lists all interactions for a customer
func (h *InteractionHandler) ListInteractionsByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	interactions, err := h.interactionService.ListByCustomer(tenantID, customerID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing interactions")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, interactions)
}

// CreateInteraction creates a new interaction
func (h *InteractionHandler) CreateInteraction(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())
	userID := auth.GetUserIDFromContext(r.Context())

	var interaction models.Interaction
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set tenant ID and created by from context
	interaction.TenantID = tenantID
	interaction.CreatedBy = userID

	// Validate interaction
	if interaction.CustomerID == "" || interaction.InteractionType == "" || interaction.InteractionDate.IsZero() {
		auth.RespondWithError(w, http.StatusBadRequest, "Customer ID, interaction type and date are required")
		return
	}

	// Check if customer exists
	customer, err := h.customerService.GetByID(tenantID, interaction.CustomerID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking customer")
		return
	}

	if customer == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	// Create interaction
	if err := h.interactionService.Create(&interaction); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating interaction")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, interaction)
}

// UpdateInteraction updates an interaction
func (h *InteractionHandler) UpdateInteraction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var interaction models.Interaction
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and tenant ID
	interaction.ID = id
	interaction.TenantID = tenantID

	// Validate interaction
	if interaction.CustomerID == "" || interaction.InteractionType == "" || interaction.InteractionDate.IsZero() {
		auth.RespondWithError(w, http.StatusBadRequest, "Customer ID, interaction type and date are required")
		return
	}

	// Check if interaction exists
	existingInteraction, err := h.interactionService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking interaction")
		return
	}

	if existingInteraction == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Interaction not found")
		return
	}

	// Update interaction
	if err := h.interactionService.Update(&interaction); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error updating interaction")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, interaction)
}

// DeleteInteraction deletes an interaction
func (h *InteractionHandler) DeleteInteraction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	// Check if interaction exists
	existingInteraction, err := h.interactionService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking interaction")
		return
	}

	if existingInteraction == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Interaction not found")
		return
	}

	// Delete interaction
	if err := h.interactionService.Delete(tenantID, id); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error deleting interaction")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Interaction deleted successfully"})
}