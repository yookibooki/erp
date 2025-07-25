package models

import (
	"time"
)

// Product represents a product in inventory
type Product struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	UnitPrice     float64   `json:"unit_price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// InventoryTransaction represents a transaction affecting inventory
type InventoryTransaction struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	ProductID       string    `json:"product_id"`
	TransactionType string    `json:"transaction_type"`
	Quantity        int       `json:"quantity"`
	Reference       string    `json:"reference"`
	Notes           string    `json:"notes"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ProductService provides methods to interact with products
type ProductService interface {
	Create(product *Product) error
	GetByID(tenantID, id string) (*Product, error)
	GetByCode(tenantID, code string) (*Product, error)
	List(tenantID string) ([]*Product, error)
	Update(product *Product) error
	Delete(tenantID, id string) error
}

// InventoryTransactionService provides methods to interact with inventory transactions
type InventoryTransactionService interface {
	Create(transaction *InventoryTransaction) error
	GetByID(tenantID, id string) (*InventoryTransaction, error)
	List(tenantID string) ([]*InventoryTransaction, error)
	ListByProduct(tenantID, productID string) ([]*InventoryTransaction, error)
}