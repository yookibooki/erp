package models

import (
	"time"
)

// Account represents a chart of account
type Account struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JournalEntry represents a journal entry
type JournalEntry struct {
	ID          string             `json:"id"`
	TenantID    string             `json:"tenant_id"`
	EntryDate   time.Time          `json:"entry_date"`
	Reference   string             `json:"reference"`
	Description string             `json:"description"`
	CreatedBy   string             `json:"created_by"`
	Lines       []JournalEntryLine `json:"lines"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// JournalEntryLine represents a line in a journal entry
type JournalEntryLine struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	JournalEntryID string    `json:"journal_entry_id"`
	AccountID      string    `json:"account_id"`
	Description    string    `json:"description"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AccountService provides methods to interact with accounts
type AccountService interface {
	Create(account *Account) error
	GetByID(tenantID, id string) (*Account, error)
	GetByCode(tenantID, code string) (*Account, error)
	List(tenantID string) ([]*Account, error)
	Update(account *Account) error
	Delete(tenantID, id string) error
}

// JournalEntryService provides methods to interact with journal entries
type JournalEntryService interface {
	Create(entry *JournalEntry) error
	GetByID(tenantID, id string) (*JournalEntry, error)
	List(tenantID string) ([]*JournalEntry, error)
	Update(entry *JournalEntry) error
	Delete(tenantID, id string) error
}