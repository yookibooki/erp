package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:12000/api"
)

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Subdomain string    `json:"subdomain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	TenantID  string `json:"tenant_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Account struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

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

type Customer struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	// Register a user
	user, token := registerUser("00000000-0000-0000-0000-000000000000")
	fmt.Printf("Registered user: %s (%s)\n", user.Email, user.ID)
	fmt.Printf("Token: %s\n", token)

	// Create a tenant
	tenant := createTenant(token)
	fmt.Printf("Created tenant: %s (%s)\n", tenant.Name, tenant.ID)

	// Create an account
	account := createAccount(token)
	fmt.Printf("Created account: %s (%s)\n", account.Name, account.ID)

	// Create a product
	product := createProduct(token)
	fmt.Printf("Created product: %s (%s)\n", product.Name, product.ID)

	// Create a customer
	customer := createCustomer(token)
	fmt.Printf("Created customer: %s (%s)\n", customer.Name, customer.ID)

	// List accounts
	accounts := listAccounts(token)
	fmt.Printf("Found %d accounts\n", len(accounts))

	// List products
	products := listProducts(token)
	fmt.Printf("Found %d products\n", len(products))

	// List customers
	customers := listCustomers(token)
	fmt.Printf("Found %d customers\n", len(customers))
}

func createTenant(token string) Tenant {
	url := fmt.Sprintf("%s/admin/tenants", baseURL)
	payload := map[string]string{
		"name":      "Test Company",
		"subdomain": "testcompany",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error creating tenant: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Error creating tenant: %s", body)
	}

	var tenant Tenant
	json.Unmarshal(body, &tenant)
	return tenant
}

func registerUser(tenantID string) (User, string) {
	url := fmt.Sprintf("%s/auth/register", baseURL)
	req := RegisterRequest{
		TenantID:  tenantID,
		Email:     "admin@testcompany.com",
		Password:  "password123",
		FirstName: "Admin",
		LastName:  "User",
		Role:      "admin",
	}

	jsonPayload, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Error registering user: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Error registering user: %s", body)
	}

	var loginResp LoginResponse
	json.Unmarshal(body, &loginResp)
	return loginResp.User, loginResp.Token
}

func createAccount(token string) Account {
	url := fmt.Sprintf("%s/accounting/accounts", baseURL)
	payload := map[string]string{
		"code":        "1000",
		"name":        "Cash",
		"type":        "asset",
		"description": "Cash on hand",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error creating account: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Error creating account: %s", body)
	}

	var account Account
	json.Unmarshal(body, &account)
	return account
}

func createProduct(token string) Product {
	url := fmt.Sprintf("%s/inventory/products", baseURL)
	payload := map[string]interface{}{
		"code":           "P001",
		"name":           "Test Product",
		"description":    "A test product",
		"unit_price":     19.99,
		"stock_quantity": 100,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error creating product: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Error creating product: %s", body)
	}

	var product Product
	json.Unmarshal(body, &product)
	return product
}

func createCustomer(token string) Customer {
	url := fmt.Sprintf("%s/crm/customers", baseURL)
	payload := map[string]string{
		"name":    "Test Customer",
		"email":   "customer@example.com",
		"phone":   "123-456-7890",
		"address": "123 Main St",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error creating customer: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Error creating customer: %s", body)
	}

	var customer Customer
	json.Unmarshal(body, &customer)
	return customer
}

func listAccounts(token string) []Account {
	url := fmt.Sprintf("%s/accounting/accounts", baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error listing accounts: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error listing accounts: %s", body)
	}

	var accounts []Account
	json.Unmarshal(body, &accounts)
	return accounts
}

func listProducts(token string) []Product {
	url := fmt.Sprintf("%s/inventory/products", baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error listing products: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error listing products: %s", body)
	}

	var products []Product
	json.Unmarshal(body, &products)
	return products
}

func listCustomers(token string) []Customer {
	url := fmt.Sprintf("%s/crm/customers", baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error listing customers: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error listing customers: %s", body)
	}

	var customers []Customer
	json.Unmarshal(body, &customers)
	return customers
}