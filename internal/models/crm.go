package models

import (
	"time"
)

// Customer represents a customer in the CRM
type Customer struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Contacts  []Contact `json:"contacts,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Contact represents a contact person for a customer
type Contact struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	CustomerID string    `json:"customer_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Position   string    `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Interaction represents an interaction with a customer
type Interaction struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	CustomerID      string    `json:"customer_id"`
	ContactID       string    `json:"contact_id,omitempty"`
	InteractionType string    `json:"interaction_type"`
	Description     string    `json:"description"`
	InteractionDate time.Time `json:"interaction_date"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CustomerService provides methods to interact with customers
type CustomerService interface {
	Create(customer *Customer) error
	GetByID(tenantID, id string) (*Customer, error)
	List(tenantID string) ([]*Customer, error)
	Update(customer *Customer) error
	Delete(tenantID, id string) error
}

// ContactService provides methods to interact with contacts
type ContactService interface {
	Create(contact *Contact) error
	GetByID(tenantID, id string) (*Contact, error)
	ListByCustomer(tenantID, customerID string) ([]*Contact, error)
	Update(contact *Contact) error
	Delete(tenantID, id string) error
}

// InteractionService provides methods to interact with interactions
type InteractionService interface {
	Create(interaction *Interaction) error
	GetByID(tenantID, id string) (*Interaction, error)
	ListByCustomer(tenantID, customerID string) ([]*Interaction, error)
	Update(interaction *Interaction) error
	Delete(tenantID, id string) error
}