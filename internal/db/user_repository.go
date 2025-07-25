package db

import (
	"database/sql"
	"time"

	"github.com/yookibooki/erp/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository implements the UserService interface
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (tenant_id, email, password_hash, first_name, last_name, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		user.TenantID,
		user.Email,
		string(hashedPassword),
		user.FirstName,
		user.LastName,
		user.Role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(tenantID, id string) (*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE tenant_id = $1 AND id = $2
	`

	user := &models.User{}
	err := r.db.QueryRow(query, tenantID, id).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(tenantID, email string) (*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE tenant_id = $1 AND email = $2
	`

	user := &models.User{}
	err := r.db.QueryRow(query, tenantID, email).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

// List lists all users for a tenant
func (r *UserRepository) List(tenantID string) ([]*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE tenant_id = $1
		ORDER BY email
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.TenantID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, first_name = $2, last_name = $3, role = $4, updated_at = $5
		WHERE tenant_id = $6 AND id = $7
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		now,
		user.TenantID,
		user.ID,
	)
	user.UpdatedAt = now
	return err
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(tenantID, id, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET password_hash = $1, updated_at = $2
		WHERE tenant_id = $3 AND id = $4
	`

	_, err = r.db.Exec(query, string(hashedPassword), time.Now(), tenantID, id)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(tenantID, id string) error {
	query := `
		DELETE FROM users
		WHERE tenant_id = $1 AND id = $2
	`

	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// Authenticate authenticates a user
func (r *UserRepository) Authenticate(tenantID, email, password string) (*models.User, error) {
	user, err := r.GetByEmail(tenantID, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}