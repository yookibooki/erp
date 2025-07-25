package db

import (
	"database/sql"
	"time"

	"github.com/yookibooki/erp/internal/models"
)

// TenantRepository implements the TenantService interface
type TenantRepository struct {
	db *DB
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create creates a new tenant
func (r *TenantRepository) Create(tenant *models.Tenant) error {
	query := `
		INSERT INTO tenants (name, subdomain)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(query, tenant.Name, tenant.Subdomain).Scan(
		&tenant.ID,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
}

// GetByID gets a tenant by ID
func (r *TenantRepository) GetByID(id string) (*models.Tenant, error) {
	query := `
		SELECT id, name, subdomain, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRow(query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return tenant, err
}

// GetBySubdomain gets a tenant by subdomain
func (r *TenantRepository) GetBySubdomain(subdomain string) (*models.Tenant, error) {
	query := `
		SELECT id, name, subdomain, created_at, updated_at
		FROM tenants
		WHERE subdomain = $1
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRow(query, subdomain).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return tenant, err
}

// List lists all tenants
func (r *TenantRepository) List() ([]*models.Tenant, error) {
	query := `
		SELECT id, name, subdomain, created_at, updated_at
		FROM tenants
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []*models.Tenant{}
	for rows.Next() {
		tenant := &models.Tenant{}
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Subdomain,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

// Update updates a tenant
func (r *TenantRepository) Update(tenant *models.Tenant) error {
	query := `
		UPDATE tenants
		SET name = $1, subdomain = $2, updated_at = $3
		WHERE id = $4
	`

	now := time.Now()
	_, err := r.db.Exec(query, tenant.Name, tenant.Subdomain, now, tenant.ID)
	tenant.UpdatedAt = now
	return err
}

// Delete deletes a tenant
func (r *TenantRepository) Delete(id string) error {
	query := `
		DELETE FROM tenants
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	return err
}