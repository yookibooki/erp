package models

import (
	"time"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Subdomain string    `json:"subdomain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TenantService provides methods to interact with tenants
type TenantService interface {
	Create(tenant *Tenant) error
	GetByID(id string) (*Tenant, error)
	GetBySubdomain(subdomain string) (*Tenant, error)
	List() ([]*Tenant, error)
	Update(tenant *Tenant) error
	Delete(id string) error
}