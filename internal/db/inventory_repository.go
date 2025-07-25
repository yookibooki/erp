package db

import (
	"database/sql"
	"time"

	"github.com/yookibooki/erp/internal/models"
)

// ProductRepository implements the ProductService interface
type ProductRepository struct {
	db *DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (tenant_id, code, name, description, unit_price, stock_quantity)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		product.TenantID,
		product.Code,
		product.Name,
		product.Description,
		product.UnitPrice,
		product.StockQuantity,
	).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
}

// GetByID gets a product by ID
func (r *ProductRepository) GetByID(tenantID, id string) (*models.Product, error) {
	query := `
		SELECT id, tenant_id, code, name, description, unit_price, stock_quantity, created_at, updated_at
		FROM products
		WHERE tenant_id = $1 AND id = $2
	`

	product := &models.Product{}
	err := r.db.QueryRow(query, tenantID, id).Scan(
		&product.ID,
		&product.TenantID,
		&product.Code,
		&product.Name,
		&product.Description,
		&product.UnitPrice,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return product, err
}

// GetByCode gets a product by code
func (r *ProductRepository) GetByCode(tenantID, code string) (*models.Product, error) {
	query := `
		SELECT id, tenant_id, code, name, description, unit_price, stock_quantity, created_at, updated_at
		FROM products
		WHERE tenant_id = $1 AND code = $2
	`

	product := &models.Product{}
	err := r.db.QueryRow(query, tenantID, code).Scan(
		&product.ID,
		&product.TenantID,
		&product.Code,
		&product.Name,
		&product.Description,
		&product.UnitPrice,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return product, err
}

// List lists all products for a tenant
func (r *ProductRepository) List(tenantID string) ([]*models.Product, error) {
	query := `
		SELECT id, tenant_id, code, name, description, unit_price, stock_quantity, created_at, updated_at
		FROM products
		WHERE tenant_id = $1
		ORDER BY code
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*models.Product{}
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.TenantID,
			&product.Code,
			&product.Name,
			&product.Description,
			&product.UnitPrice,
			&product.StockQuantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// Update updates a product
func (r *ProductRepository) Update(product *models.Product) error {
	query := `
		UPDATE products
		SET code = $1, name = $2, description = $3, unit_price = $4, stock_quantity = $5, updated_at = $6
		WHERE tenant_id = $7 AND id = $8
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		product.Code,
		product.Name,
		product.Description,
		product.UnitPrice,
		product.StockQuantity,
		now,
		product.TenantID,
		product.ID,
	)
	product.UpdatedAt = now
	return err
}

// Delete deletes a product
func (r *ProductRepository) Delete(tenantID, id string) error {
	query := `
		DELETE FROM products
		WHERE tenant_id = $1 AND id = $2
	`

	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// InventoryTransactionRepository implements the InventoryTransactionService interface
type InventoryTransactionRepository struct {
	db *DB
}

// NewInventoryTransactionRepository creates a new inventory transaction repository
func NewInventoryTransactionRepository(db *DB) *InventoryTransactionRepository {
	return &InventoryTransactionRepository{db: db}
}

// Create creates a new inventory transaction
func (r *InventoryTransactionRepository) Create(transaction *models.InventoryTransaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert inventory transaction
	query := `
		INSERT INTO inventory_transactions (tenant_id, product_id, transaction_type, quantity, reference, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(
		query,
		transaction.TenantID,
		transaction.ProductID,
		transaction.TransactionType,
		transaction.Quantity,
		transaction.Reference,
		transaction.Notes,
		transaction.CreatedBy,
	).Scan(
		&transaction.ID,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Update product stock quantity
	var stockChange int
	if transaction.TransactionType == "IN" {
		stockChange = transaction.Quantity
	} else if transaction.TransactionType == "OUT" {
		stockChange = -transaction.Quantity
	}

	if stockChange != 0 {
		query := `
			UPDATE products
			SET stock_quantity = stock_quantity + $1, updated_at = $2
			WHERE tenant_id = $3 AND id = $4
		`

		_, err = tx.Exec(
			query,
			stockChange,
			time.Now(),
			transaction.TenantID,
			transaction.ProductID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByID gets an inventory transaction by ID
func (r *InventoryTransactionRepository) GetByID(tenantID, id string) (*models.InventoryTransaction, error) {
	query := `
		SELECT id, tenant_id, product_id, transaction_type, quantity, reference, notes, created_by, created_at, updated_at
		FROM inventory_transactions
		WHERE tenant_id = $1 AND id = $2
	`

	transaction := &models.InventoryTransaction{}
	err := r.db.QueryRow(query, tenantID, id).Scan(
		&transaction.ID,
		&transaction.TenantID,
		&transaction.ProductID,
		&transaction.TransactionType,
		&transaction.Quantity,
		&transaction.Reference,
		&transaction.Notes,
		&transaction.CreatedBy,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return transaction, err
}

// List lists all inventory transactions for a tenant
func (r *InventoryTransactionRepository) List(tenantID string) ([]*models.InventoryTransaction, error) {
	query := `
		SELECT id, tenant_id, product_id, transaction_type, quantity, reference, notes, created_by, created_at, updated_at
		FROM inventory_transactions
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []*models.InventoryTransaction{}
	for rows.Next() {
		transaction := &models.InventoryTransaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.TenantID,
			&transaction.ProductID,
			&transaction.TransactionType,
			&transaction.Quantity,
			&transaction.Reference,
			&transaction.Notes,
			&transaction.CreatedBy,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// ListByProduct lists all inventory transactions for a product
func (r *InventoryTransactionRepository) ListByProduct(tenantID, productID string) ([]*models.InventoryTransaction, error) {
	query := `
		SELECT id, tenant_id, product_id, transaction_type, quantity, reference, notes, created_by, created_at, updated_at
		FROM inventory_transactions
		WHERE tenant_id = $1 AND product_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, tenantID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []*models.InventoryTransaction{}
	for rows.Next() {
		transaction := &models.InventoryTransaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.TenantID,
			&transaction.ProductID,
			&transaction.TransactionType,
			&transaction.Quantity,
			&transaction.Reference,
			&transaction.Notes,
			&transaction.CreatedBy,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}