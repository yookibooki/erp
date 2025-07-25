package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserService provides methods to interact with users
type UserService interface {
	Create(user *User) error
	GetByID(tenantID, id string) (*User, error)
	GetByEmail(tenantID, email string) (*User, error)
	List(tenantID string) ([]*User, error)
	Update(user *User) error
	Delete(tenantID, id string) error
	Authenticate(tenantID, email, password string) (*User, error)
}