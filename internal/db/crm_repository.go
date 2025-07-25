package db

import (
	"database/sql"
	"time"

	"github.com/yookibooki/erp/internal/models"
)

// CustomerRepository implements the CustomerService interface
type CustomerRepository struct {
	db *DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create creates a new customer
func (r *CustomerRepository) Create(customer *models.Customer) error {
	query := `
		INSERT INTO customers (tenant_id, name, email, phone, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		customer.TenantID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
	).Scan(
		&customer.ID,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
}

// GetByID gets a customer by ID
func (r *CustomerRepository) GetByID(tenantID, id string) (*models.Customer, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, address, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1 AND id = $2
	`

	customer := &models.Customer{}
	err := r.db.QueryRow(query, tenantID, id).Scan(
		&customer.ID,
		&customer.TenantID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return customer, err
}

// List lists all customers for a tenant
func (r *CustomerRepository) List(tenantID string) ([]*models.Customer, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, address, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		ORDER BY name
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []*models.Customer{}
	for rows.Next() {
		customer := &models.Customer{}
		err := rows.Scan(
			&customer.ID,
			&customer.TenantID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

// Update updates a customer
func (r *CustomerRepository) Update(customer *models.Customer) error {
	query := `
		UPDATE customers
		SET name = $1, email = $2, phone = $3, address = $4, updated_at = $5
		WHERE tenant_id = $6 AND id = $7
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		now,
		customer.TenantID,
		customer.ID,
	)
	customer.UpdatedAt = now
	return err
}

// Delete deletes a customer
func (r *CustomerRepository) Delete(tenantID, id string) error {
	query := `
		DELETE FROM customers
		WHERE tenant_id = $1 AND id = $2
	`

	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// ContactRepository implements the ContactService interface
type ContactRepository struct {
	db *DB
}

// NewContactRepository creates a new contact repository
func NewContactRepository(db *DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// Create creates a new contact
func (r *ContactRepository) Create(contact *models.Contact) error {
	query := `
		INSERT INTO contacts (tenant_id, customer_id, first_name, last_name, email, phone, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		contact.TenantID,
		contact.CustomerID,
		contact.FirstName,
		contact.LastName,
		contact.Email,
		contact.Phone,
		contact.Position,
	).Scan(
		&contact.ID,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)
}

// GetByID gets a contact by ID
func (r *ContactRepository) GetByID(tenantID, id string) (*models.Contact, error) {
	query := `
		SELECT id, tenant_id, customer_id, first_name, last_name, email, phone, position, created_at, updated_at
		FROM contacts
		WHERE tenant_id = $1 AND id = $2
	`

	contact := &models.Contact{}
	err := r.db.QueryRow(query, tenantID, id).Scan(
		&contact.ID,
		&contact.TenantID,
		&contact.CustomerID,
		&contact.FirstName,
		&contact.LastName,
		&contact.Email,
		&contact.Phone,
		&contact.Position,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return contact, err
}

// ListByCustomer lists all contacts for a customer
func (r *ContactRepository) ListByCustomer(tenantID, customerID string) ([]*models.Contact, error) {
	query := `
		SELECT id, tenant_id, customer_id, first_name, last_name, email, phone, position, created_at, updated_at
		FROM contacts
		WHERE tenant_id = $1 AND customer_id = $2
		ORDER BY last_name, first_name
	`

	rows, err := r.db.Query(query, tenantID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := []*models.Contact{}
	for rows.Next() {
		contact := &models.Contact{}
		err := rows.Scan(
			&contact.ID,
			&contact.TenantID,
			&contact.CustomerID,
			&contact.FirstName,
			&contact.LastName,
			&contact.Email,
			&contact.Phone,
			&contact.Position,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// Update updates a contact
func (r *ContactRepository) Update(contact *models.Contact) error {
	query := `
		UPDATE contacts
		SET customer_id = $1, first_name = $2, last_name = $3, email = $4, phone = $5, position = $6, updated_at = $7
		WHERE tenant_id = $8 AND id = $9
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		contact.CustomerID,
		contact.FirstName,
		contact.LastName,
		contact.Email,
		contact.Phone,
		contact.Position,
		now,
		contact.TenantID,
		contact.ID,
	)
	contact.UpdatedAt = now
	return err
}

// Delete deletes a contact
func (r *ContactRepository) Delete(tenantID, id string) error {
	query := `
		DELETE FROM contacts
		WHERE tenant_id = $1 AND id = $2
	`

	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// InteractionRepository implements the InteractionService interface
type InteractionRepository struct {
	db *DB
}

// NewInteractionRepository creates a new interaction repository
func NewInteractionRepository(db *DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

// Create creates a new interaction
func (r *InteractionRepository) Create(interaction *models.Interaction) error {
	query := `
		INSERT INTO interactions (tenant_id, customer_id, contact_id, interaction_type, description, interaction_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		interaction.TenantID,
		interaction.CustomerID,
		interaction.ContactID,
		interaction.InteractionType,
		interaction.Description,
		interaction.InteractionDate,
		interaction.CreatedBy,
	).Scan(
		&interaction.ID,
		&interaction.CreatedAt,
		&interaction.UpdatedAt,
	)
}

// GetByID gets an interaction by ID
func (r *InteractionRepository) GetByID(tenantID, id string) (*models.Interaction, error) {
	query := `
		SELECT id, tenant_id, customer_id, contact_id, interaction_type, description, interaction_date, created_by, created_at, updated_at
		FROM interactions
		WHERE tenant_id = $1 AND id = $2
	`

	interaction := &models.Interaction{}
	var contactID sql.NullString
	err := r.db.QueryRow(query, tenantID, id).Scan(
		&interaction.ID,
		&interaction.TenantID,
		&interaction.CustomerID,
		&contactID,
		&interaction.InteractionType,
		&interaction.Description,
		&interaction.InteractionDate,
		&interaction.CreatedBy,
		&interaction.CreatedAt,
		&interaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if contactID.Valid {
		interaction.ContactID = contactID.String
	}

	return interaction, nil
}

// ListByCustomer lists all interactions for a customer
func (r *InteractionRepository) ListByCustomer(tenantID, customerID string) ([]*models.Interaction, error) {
	query := `
		SELECT id, tenant_id, customer_id, contact_id, interaction_type, description, interaction_date, created_by, created_at, updated_at
		FROM interactions
		WHERE tenant_id = $1 AND customer_id = $2
		ORDER BY interaction_date DESC
	`

	rows, err := r.db.Query(query, tenantID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	interactions := []*models.Interaction{}
	for rows.Next() {
		interaction := &models.Interaction{}
		var contactID sql.NullString
		err := rows.Scan(
			&interaction.ID,
			&interaction.TenantID,
			&interaction.CustomerID,
			&contactID,
			&interaction.InteractionType,
			&interaction.Description,
			&interaction.InteractionDate,
			&interaction.CreatedBy,
			&interaction.CreatedAt,
			&interaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if contactID.Valid {
			interaction.ContactID = contactID.String
		}

		interactions = append(interactions, interaction)
	}

	return interactions, nil
}

// Update updates an interaction
func (r *InteractionRepository) Update(interaction *models.Interaction) error {
	query := `
		UPDATE interactions
		SET customer_id = $1, contact_id = $2, interaction_type = $3, description = $4, interaction_date = $5, updated_at = $6
		WHERE tenant_id = $7 AND id = $8
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		interaction.CustomerID,
		interaction.ContactID,
		interaction.InteractionType,
		interaction.Description,
		interaction.InteractionDate,
		now,
		interaction.TenantID,
		interaction.ID,
	)
	interaction.UpdatedAt = now
	return err
}

// Delete deletes an interaction
func (r *InteractionRepository) Delete(tenantID, id string) error {
	query := `
		DELETE FROM interactions
		WHERE tenant_id = $1 AND id = $2
	`

	_, err := r.db.Exec(query, tenantID, id)
	return err
}