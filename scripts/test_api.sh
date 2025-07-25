#!/bin/bash

# Test API script for ERP SaaS

BASE_URL="http://localhost:12000/api"
TOKEN=""

# First, register a user with a new tenant
echo "Registering user with new tenant..."
USER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "00000000-0000-0000-0000-000000000000",
    "email": "admin@testcompany.com",
    "password": "password123",
    "first_name": "Admin",
    "last_name": "User",
    "role": "admin"
  }')
echo $USER_RESPONSE
TOKEN=$(echo $USER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token: $TOKEN"

# Create a tenant
echo "Creating tenant..."
TENANT_RESPONSE=$(curl -s -X POST "${BASE_URL}/admin/tenants" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Company",
    "subdomain": "testcompany"
  }')
echo $TENANT_RESPONSE
TENANT_ID=$(echo $TENANT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Tenant ID: $TENANT_ID"

# Register a user
echo "Registering user..."
USER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"tenant_id\": \"$TENANT_ID\",
    \"email\": \"admin@testcompany.com\",
    \"password\": \"password123\",
    \"first_name\": \"Admin\",
    \"last_name\": \"User\",
    \"role\": \"admin\"
  }")
echo $USER_RESPONSE
TOKEN=$(echo $USER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token: $TOKEN"

# Create an account
echo "Creating account..."
ACCOUNT_RESPONSE=$(curl -s -X POST "${BASE_URL}/accounting/accounts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "code": "1000",
    "name": "Cash",
    "type": "asset",
    "description": "Cash on hand"
  }')
echo $ACCOUNT_RESPONSE
ACCOUNT_ID=$(echo $ACCOUNT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Account ID: $ACCOUNT_ID"

# Create a product
echo "Creating product..."
PRODUCT_RESPONSE=$(curl -s -X POST "${BASE_URL}/inventory/products" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "code": "P001",
    "name": "Test Product",
    "description": "A test product",
    "unit_price": 19.99,
    "stock_quantity": 100
  }')
echo $PRODUCT_RESPONSE
PRODUCT_ID=$(echo $PRODUCT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Product ID: $PRODUCT_ID"

# Create a customer
echo "Creating customer..."
CUSTOMER_RESPONSE=$(curl -s -X POST "${BASE_URL}/crm/customers" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Customer",
    "email": "customer@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St"
  }')
echo $CUSTOMER_RESPONSE
CUSTOMER_ID=$(echo $CUSTOMER_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Customer ID: $CUSTOMER_ID"

# Create a contact
echo "Creating contact..."
CONTACT_RESPONSE=$(curl -s -X POST "${BASE_URL}/crm/contacts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"customer_id\": \"$CUSTOMER_ID\",
    \"first_name\": \"John\",
    \"last_name\": \"Doe\",
    \"email\": \"john.doe@example.com\",
    \"phone\": \"123-456-7890\",
    \"position\": \"CEO\"
  }")
echo $CONTACT_RESPONSE

# List accounts
echo "Listing accounts..."
curl -s -X GET "${BASE_URL}/accounting/accounts" \
  -H "Authorization: Bearer $TOKEN"

# List products
echo "Listing products..."
curl -s -X GET "${BASE_URL}/inventory/products" \
  -H "Authorization: Bearer $TOKEN"

# List customers
echo "Listing customers..."
curl -s -X GET "${BASE_URL}/crm/customers" \
  -H "Authorization: Bearer $TOKEN"