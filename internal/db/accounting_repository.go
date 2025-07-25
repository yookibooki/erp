package db

import (
	"database/sql"
	"time"

	"github.com/yookibooki/erp/internal/models"
)

// AccountRepository implements the AccountService interface
type AccountRepository struct {
	db *DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account
func (r *AccountRepository) Create(account *models.Account) error {
	query := `
		INSERT INTO accounts (tenant_id, code, name, type, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		account.TenantID,
		account.Code,
		account.Name,
		account.Type,
		account.Description,
	).Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
}

// GetByID gets an account by ID
func (r *AccountRepository) GetByID(tenantID, id string) (*models.Account, error) {
	query := `
		SELECT id, tenant_id, code, name, type, description, created_at, updated_at
		FROM accounts
		WHERE tenant_id = $1 AND id = $2
	`

	account := &models.Account{}
	err := r.db.QueryRow(query, tenantID, id).Scan(
		&account.ID,
		&account.TenantID,
		&account.Code,
		&account.Name,
		&account.Type,
		&account.Description,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return account, err
}

// GetByCode gets an account by code
func (r *AccountRepository) GetByCode(tenantID, code string) (*models.Account, error) {
	query := `
		SELECT id, tenant_id, code, name, type, description, created_at, updated_at
		FROM accounts
		WHERE tenant_id = $1 AND code = $2
	`

	account := &models.Account{}
	err := r.db.QueryRow(query, tenantID, code).Scan(
		&account.ID,
		&account.TenantID,
		&account.Code,
		&account.Name,
		&account.Type,
		&account.Description,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return account, err
}

// List lists all accounts for a tenant
func (r *AccountRepository) List(tenantID string) ([]*models.Account, error) {
	query := `
		SELECT id, tenant_id, code, name, type, description, created_at, updated_at
		FROM accounts
		WHERE tenant_id = $1
		ORDER BY code
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*models.Account{}
	for rows.Next() {
		account := &models.Account{}
		err := rows.Scan(
			&account.ID,
			&account.TenantID,
			&account.Code,
			&account.Name,
			&account.Type,
			&account.Description,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// Update updates an account
func (r *AccountRepository) Update(account *models.Account) error {
	query := `
		UPDATE accounts
		SET code = $1, name = $2, type = $3, description = $4, updated_at = $5
		WHERE tenant_id = $6 AND id = $7
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		account.Code,
		account.Name,
		account.Type,
		account.Description,
		now,
		account.TenantID,
		account.ID,
	)
	account.UpdatedAt = now
	return err
}

// Delete deletes an account
func (r *AccountRepository) Delete(tenantID, id string) error {
	query := `
		DELETE FROM accounts
		WHERE tenant_id = $1 AND id = $2
	`

	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// JournalEntryRepository implements the JournalEntryService interface
type JournalEntryRepository struct {
	db *DB
}

// NewJournalEntryRepository creates a new journal entry repository
func NewJournalEntryRepository(db *DB) *JournalEntryRepository {
	return &JournalEntryRepository{db: db}
}

// Create creates a new journal entry
func (r *JournalEntryRepository) Create(entry *models.JournalEntry) error {
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

	// Insert journal entry
	query := `
		INSERT INTO journal_entries (tenant_id, entry_date, reference, description, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(
		query,
		entry.TenantID,
		entry.EntryDate,
		entry.Reference,
		entry.Description,
		entry.CreatedBy,
	).Scan(
		&entry.ID,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert journal entry lines
	for i := range entry.Lines {
		line := &entry.Lines[i]
		line.JournalEntryID = entry.ID

		query := `
			INSERT INTO journal_entry_lines (tenant_id, journal_entry_id, account_id, description, debit, credit)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at
		`

		err = tx.QueryRow(
			query,
			entry.TenantID,
			line.JournalEntryID,
			line.AccountID,
			line.Description,
			line.Debit,
			line.Credit,
		).Scan(
			&line.ID,
			&line.CreatedAt,
			&line.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByID gets a journal entry by ID
func (r *JournalEntryRepository) GetByID(tenantID, id string) (*models.JournalEntry, error) {
	// Get journal entry
	query := `
		SELECT id, tenant_id, entry_date, reference, description, created_by, created_at, updated_at
		FROM journal_entries
		WHERE tenant_id = $1 AND id = $2
	`

	entry := &models.JournalEntry{}
	err := r.db.QueryRow(query, tenantID, id).Scan(
		&entry.ID,
		&entry.TenantID,
		&entry.EntryDate,
		&entry.Reference,
		&entry.Description,
		&entry.CreatedBy,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get journal entry lines
	query = `
		SELECT id, tenant_id, journal_entry_id, account_id, description, debit, credit, created_at, updated_at
		FROM journal_entry_lines
		WHERE tenant_id = $1 AND journal_entry_id = $2
		ORDER BY id
	`

	rows, err := r.db.Query(query, tenantID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entry.Lines = []models.JournalEntryLine{}
	for rows.Next() {
		line := models.JournalEntryLine{}
		err := rows.Scan(
			&line.ID,
			&line.TenantID,
			&line.JournalEntryID,
			&line.AccountID,
			&line.Description,
			&line.Debit,
			&line.Credit,
			&line.CreatedAt,
			&line.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entry.Lines = append(entry.Lines, line)
	}

	return entry, nil
}

// List lists all journal entries for a tenant
func (r *JournalEntryRepository) List(tenantID string) ([]*models.JournalEntry, error) {
	query := `
		SELECT id, tenant_id, entry_date, reference, description, created_by, created_at, updated_at
		FROM journal_entries
		WHERE tenant_id = $1
		ORDER BY entry_date DESC
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []*models.JournalEntry{}
	for rows.Next() {
		entry := &models.JournalEntry{}
		err := rows.Scan(
			&entry.ID,
			&entry.TenantID,
			&entry.EntryDate,
			&entry.Reference,
			&entry.Description,
			&entry.CreatedBy,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	// Get journal entry lines for each entry
	for _, entry := range entries {
		query := `
			SELECT id, tenant_id, journal_entry_id, account_id, description, debit, credit, created_at, updated_at
			FROM journal_entry_lines
			WHERE tenant_id = $1 AND journal_entry_id = $2
			ORDER BY id
		`

		rows, err := r.db.Query(query, tenantID, entry.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		entry.Lines = []models.JournalEntryLine{}
		for rows.Next() {
			line := models.JournalEntryLine{}
			err := rows.Scan(
				&line.ID,
				&line.TenantID,
				&line.JournalEntryID,
				&line.AccountID,
				&line.Description,
				&line.Debit,
				&line.Credit,
				&line.CreatedAt,
				&line.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
			entry.Lines = append(entry.Lines, line)
		}
	}

	return entries, nil
}

// Update updates a journal entry
func (r *JournalEntryRepository) Update(entry *models.JournalEntry) error {
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

	// Update journal entry
	query := `
		UPDATE journal_entries
		SET entry_date = $1, reference = $2, description = $3, updated_at = $4
		WHERE tenant_id = $5 AND id = $6
	`

	now := time.Now()
	_, err = tx.Exec(
		query,
		entry.EntryDate,
		entry.Reference,
		entry.Description,
		now,
		entry.TenantID,
		entry.ID,
	)
	if err != nil {
		return err
	}
	entry.UpdatedAt = now

	// Delete existing journal entry lines
	query = `
		DELETE FROM journal_entry_lines
		WHERE tenant_id = $1 AND journal_entry_id = $2
	`

	_, err = tx.Exec(query, entry.TenantID, entry.ID)
	if err != nil {
		return err
	}

	// Insert new journal entry lines
	for i := range entry.Lines {
		line := &entry.Lines[i]
		line.JournalEntryID = entry.ID

		query := `
			INSERT INTO journal_entry_lines (tenant_id, journal_entry_id, account_id, description, debit, credit)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at
		`

		err = tx.QueryRow(
			query,
			entry.TenantID,
			line.JournalEntryID,
			line.AccountID,
			line.Description,
			line.Debit,
			line.Credit,
		).Scan(
			&line.ID,
			&line.CreatedAt,
			&line.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a journal entry
func (r *JournalEntryRepository) Delete(tenantID, id string) error {
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

	// Delete journal entry lines
	query := `
		DELETE FROM journal_entry_lines
		WHERE tenant_id = $1 AND journal_entry_id = $2
	`

	_, err = tx.Exec(query, tenantID, id)
	if err != nil {
		return err
	}

	// Delete journal entry
	query = `
		DELETE FROM journal_entries
		WHERE tenant_id = $1 AND id = $2
	`

	_, err = tx.Exec(query, tenantID, id)
	return err
}