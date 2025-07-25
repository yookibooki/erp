package inventory

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yookibooki/erp/internal/auth"
	"github.com/yookibooki/erp/internal/models"
)

// ProductHandler handles product requests
type ProductHandler struct {
	productService models.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService models.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProduct gets a product by ID
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	product, err := h.productService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting product")
		return
	}

	if product == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, product)
}

// ListProducts lists all products for a tenant
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	products, err := h.productService.List(tenantID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing products")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, products)
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set tenant ID from context
	product.TenantID = tenantID

	// Validate product
	if product.Code == "" || product.Name == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Code and name are required")
		return
	}

	// Check if product already exists
	existingProduct, err := h.productService.GetByCode(tenantID, product.Code)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking product")
		return
	}

	if existingProduct != nil {
		auth.RespondWithError(w, http.StatusConflict, "Product with this code already exists")
		return
	}

	// Create product
	if err := h.productService.Create(&product); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating product")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, product)
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and tenant ID
	product.ID = id
	product.TenantID = tenantID

	// Validate product
	if product.Code == "" || product.Name == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Code and name are required")
		return
	}

	// Check if product exists
	existingProduct, err := h.productService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking product")
		return
	}

	if existingProduct == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Update product
	if err := h.productService.Update(&product); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error updating product")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, product)
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	// Check if product exists
	existingProduct, err := h.productService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking product")
		return
	}

	if existingProduct == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Delete product
	if err := h.productService.Delete(tenantID, id); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error deleting product")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

// InventoryTransactionHandler handles inventory transaction requests
type InventoryTransactionHandler struct {
	transactionService models.InventoryTransactionService
	productService     models.ProductService
}

// NewInventoryTransactionHandler creates a new inventory transaction handler
func NewInventoryTransactionHandler(
	transactionService models.InventoryTransactionService,
	productService models.ProductService,
) *InventoryTransactionHandler {
	return &InventoryTransactionHandler{
		transactionService: transactionService,
		productService:     productService,
	}
}

// GetTransaction gets a transaction by ID
func (h *InventoryTransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	transaction, err := h.transactionService.GetByID(tenantID, id)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error getting transaction")
		return
	}

	if transaction == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Transaction not found")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, transaction)
}

// ListTransactions lists all transactions for a tenant
func (h *InventoryTransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())

	transactions, err := h.transactionService.List(tenantID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing transactions")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, transactions)
}

// ListTransactionsByProduct lists all transactions for a product
func (h *InventoryTransactionHandler) ListTransactionsByProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]
	tenantID := auth.GetTenantIDFromContext(r.Context())

	transactions, err := h.transactionService.ListByProduct(tenantID, productID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error listing transactions")
		return
	}

	auth.RespondWithJSON(w, http.StatusOK, transactions)
}

// CreateTransaction creates a new transaction
func (h *InventoryTransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	tenantID := auth.GetTenantIDFromContext(r.Context())
	userID := auth.GetUserIDFromContext(r.Context())

	var transaction models.InventoryTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		auth.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set tenant ID and created by from context
	transaction.TenantID = tenantID
	transaction.CreatedBy = userID

	// Validate transaction
	if transaction.ProductID == "" || transaction.TransactionType == "" {
		auth.RespondWithError(w, http.StatusBadRequest, "Product ID and transaction type are required")
		return
	}

	// Check if product exists
	product, err := h.productService.GetByID(tenantID, transaction.ProductID)
	if err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error checking product")
		return
	}

	if product == nil {
		auth.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Create transaction
	if err := h.transactionService.Create(&transaction); err != nil {
		auth.RespondWithError(w, http.StatusInternalServerError, "Error creating transaction")
		return
	}

	auth.RespondWithJSON(w, http.StatusCreated, transaction)
}