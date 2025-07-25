package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yookibooki/erp/internal/auth"
	"github.com/yookibooki/erp/internal/models"
)

// AccountHandler handles account requests
type AccountHandler struct {
	accountService models.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService models.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// GetAccount gets an account by ID
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	account, err := h.accountService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting account")
		return
	}

	if account == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Account not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, account)
}

// ListAccounts lists all accounts for a tenant
func (h *AccountHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	accounts, err := h.accountService.List(tenantID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing accounts")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, accounts)
}

// CreateAccount creates a new account
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set tenant ID from context
	account.TenantID = tenantID

	// Validate account
	if account.Code == "" || account.Name == "" || account.Type == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Code, name and type are required")
		return
	}

	// Check if account already exists
	existingAccount, err := h.accountService.GetByCode(tenantID, account.Code)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking account")
		return
	}

	if existingAccount != nil {
		auth.RespondWithError(w, http.StatusConflict, "Account with this code already exists")
		return
	}

	// Create account
	if err := h.accountService.Create(&account); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating account")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, account)
}

// UpdateAccount updates an account
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and tenant ID
	account.ID = id
	account.TenantID = tenantID

	// Validate account
	if account.Code == "" || account.Name == "" || account.Type == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Code, name and type are required")
		return
	}

	// Check if account exists
	existingAccount, err := h.accountService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking account")
		return
	}

	if existingAccount == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Account not found")
		return
	}

	// Update account
	if err := h.accountService.Update(&account); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error updating account")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, account)
}

// DeleteAccount deletes an account
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	// Check if account exists
	existingAccount, err := h.accountService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking account")
		return
	}

	if existingAccount == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Account not found")
		return
	}

	// Delete account
	if err := h.accountService.Delete(tenantID, id); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error deleting account")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}

// JournalEntryHandler handles journal entry requests
type JournalEntryHandler struct {
	journalEntryService models.JournalEntryService
}

// NewJournalEntryHandler creates a new journal entry handler
func NewJournalEntryHandler(journalEntryService models.JournalEntryService) *JournalEntryHandler {
	return &JournalEntryHandler{
		journalEntryService: journalEntryService,
	}
}

// GetJournalEntry gets a journal entry by ID
func (h *JournalEntryHandler) GetJournalEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	entry, err := h.journalEntryService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting journal entry")
		return
	}

	if entry == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Journal entry not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, entry)
}

// ListJournalEntries lists all journal entries for a tenant
func (h *JournalEntryHandler) ListJournalEntries(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	entries, err := h.journalEntryService.List(tenantID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing journal entries")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, entries)
}

// CreateJournalEntry creates a new journal entry
func (h *JournalEntryHandler) CreateJournalEntry(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())
	userID := auth.GetUserIDFromContext(r.Context())

	var entry models.JournalEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set tenant ID and created by from context
	entry.TenantID = tenantID
	entry.CreatedBy = userID

	// Validate entry
	if entry.EntryDate.IsZero() {
		auth.RespondWithError(w, http.StatusBadRequest, "Entry date is required")
		return
	}

	if len(entry.Lines) == 0 {
		auth.RespondWithError(w, http.StatusBadRequest, "At least one journal entry line is required")
		return
	}

	// Set tenant ID for all lines
	for i := range entry.Lines {
		entry.Lines[i].TenantID = tenantID
	}

	// Create journal entry
	if err := h.journalEntryService.Create(&entry); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating journal entry")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, entry)
}

// UpdateJournalEntry updates a journal entry
func (h *JournalEntryHandler) UpdateJournalEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var entry models.JournalEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and tenant ID
	entry.ID = id
	entry.TenantID = tenantID

	// Validate entry
	if entry.EntryDate.IsZero() {
		auth.RespondWithError(w, http.StatusBadRequest, "Entry date is required")
		return
	}

	if len(entry.Lines) == 0 {
		auth.RespondWithError(w, http.StatusBadRequest, "At least one journal entry line is required")
		return
	}

	// Set tenant ID for all lines
	for i := range entry.Lines {
		entry.Lines[i].TenantID = tenantID
		entry.Lines[i].JournalEntryID = id
	}

	// Check if journal entry exists
	existingEntry, err := h.journalEntryService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking journal entry")
		return
	}

	if existingEntry == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Journal entry not found")
		return
	}

	// Update journal entry
	if err := h.journalEntryService.Update(&entry); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error updating journal entry")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, entry)
}

// DeleteJournalEntry deletes a journal entry
func (h *JournalEntryHandler) DeleteJournalEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	// Check if journal entry exists
	existingEntry, err := h.journalEntryService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking journal entry")
		return
	}

	if existingEntry == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Journal entry not found")
		return
	}

	// Delete journal entry
	if err := h.journalEntryService.Delete(tenantID, id); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error deleting journal entry")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Journal entry deleted successfully"})
}