# Multi-Tenant SaaS ERP

A minimalist, multi-tenant SaaS ERP system built with Go and PostgreSQL.

## Features

- **Multi-tenant Architecture**: Uses a shared database with tenant_id for data isolation
- **Authentication**: JWT-based authentication and authorization
- **Core Modules**:
  - **Accounting**: Chart of accounts, journal entries
  - **Inventory**: Products, inventory transactions
  - **CRM**: Customers, contacts, interactions

## Tech Stack

- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Authentication**: JWT
- **API**: RESTful JSON API

## Project Structure

```
.
├── cmd
│   └── api                 # Application entry point
├── internal
│   ├── api                 # API handlers
│   ├── auth                # Authentication
│   ├── config              # Configuration
│   ├── db                  # Database connection and repositories
│   ├── models              # Data models
│   └── modules             # Business modules
│       ├── accounting      # Accounting module
│       ├── inventory       # Inventory module
│       └── crm             # CRM module
├── migrations              # Database migrations
└── scripts                 # Utility scripts
```

## Getting Started

### Prerequisites

- Go 1.19+
- PostgreSQL 15+

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/yookibooki/erp.git
   cd erp
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up the database:
   ```
   # Create database and user
   sudo -u postgres psql -c "CREATE DATABASE erp_saas;"
   sudo -u postgres psql -c "CREATE USER erp_user WITH PASSWORD 'erp_password';"
   sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE erp_saas TO erp_user;"
   sudo -u postgres psql -c "GRANT ALL ON SCHEMA public TO erp_user;" -d erp_saas
   
   # Run migrations
   cd scripts
   ./run_migrations.sh
   ```

4. Run the application:
   ```
   go run cmd/api/main.go
   ```

## API Endpoints

### Authentication

- `POST /api/auth/login`: Login
- `POST /api/auth/register`: Register

### Tenants

- `GET /api/tenants/{subdomain}`: Get tenant by subdomain
- `GET /api/admin/tenants`: List all tenants
- `POST /api/admin/tenants`: Create a new tenant
- `GET /api/admin/tenants/{id}`: Get tenant by ID
- `PUT /api/admin/tenants/{id}`: Update tenant
- `DELETE /api/admin/tenants/{id}`: Delete tenant

### Users

- `GET /api/users`: List all users
- `POST /api/users`: Create a new user
- `GET /api/users/{id}`: Get user by ID
- `PUT /api/users/{id}`: Update user
- `DELETE /api/users/{id}`: Delete user

### Accounting

- `GET /api/accounting/accounts`: List all accounts
- `POST /api/accounting/accounts`: Create a new account
- `GET /api/accounting/accounts/{id}`: Get account by ID
- `PUT /api/accounting/accounts/{id}`: Update account
- `DELETE /api/accounting/accounts/{id}`: Delete account

- `GET /api/accounting/journal-entries`: List all journal entries
- `POST /api/accounting/journal-entries`: Create a new journal entry
- `GET /api/accounting/journal-entries/{id}`: Get journal entry by ID
- `PUT /api/accounting/journal-entries/{id}`: Update journal entry
- `DELETE /api/accounting/journal-entries/{id}`: Delete journal entry

### Inventory

- `GET /api/inventory/products`: List all products
- `POST /api/inventory/products`: Create a new product
- `GET /api/inventory/products/{id}`: Get product by ID
- `PUT /api/inventory/products/{id}`: Update product
- `DELETE /api/inventory/products/{id}`: Delete product

- `GET /api/inventory/transactions`: List all inventory transactions
- `POST /api/inventory/transactions`: Create a new inventory transaction
- `GET /api/inventory/transactions/{id}`: Get inventory transaction by ID
- `GET /api/inventory/transactions/product/{productId}`: List transactions by product

### CRM

- `GET /api/crm/customers`: List all customers
- `POST /api/crm/customers`: Create a new customer
- `GET /api/crm/customers/{id}`: Get customer by ID
- `PUT /api/crm/customers/{id}`: Update customer
- `DELETE /api/crm/customers/{id}`: Delete customer

- `POST /api/crm/contacts`: Create a new contact
- `GET /api/crm/contacts/{id}`: Get contact by ID
- `PUT /api/crm/contacts/{id}`: Update contact
- `DELETE /api/crm/contacts/{id}`: Delete contact
- `GET /api/crm/customers/{customerId}/contacts`: List contacts by customer

- `POST /api/crm/interactions`: Create a new interaction
- `GET /api/crm/interactions/{id}`: Get interaction by ID
- `PUT /api/crm/interactions/{id}`: Update interaction
- `DELETE /api/crm/interactions/{id}`: Delete interaction
- `GET /api/crm/customers/{customerId}/interactions`: List interactions by customer

## License

This project is licensed under the MIT License - see the LICENSE file for details.