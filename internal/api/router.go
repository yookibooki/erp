package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yookibooki/erp/internal/auth"
	"github.com/yookibooki/erp/internal/models"
	"github.com/yookibooki/erp/internal/modules/accounting"
	"github.com/yookibooki/erp/internal/modules/crm"
	"github.com/yookibooki/erp/internal/modules/inventory"
)

// Router is the HTTP router
type Router struct {
	*mux.Router
	jwtService *auth.JWTService
}

// NewRouter creates a new router
func NewRouter(
	tenantService models.TenantService,
	userService models.UserService,
	accountService models.AccountService,
	journalEntryService models.JournalEntryService,
	productService models.ProductService,
	inventoryTransactionService models.InventoryTransactionService,
	customerService models.CustomerService,
	contactService models.ContactService,
	interactionService models.InteractionService,
	jwtService *auth.JWTService,
) *Router {
	r := mux.NewRouter()
	router := &Router{r, jwtService}

	// Create handlers
	authHandler := NewAuthHandler(userService, jwtService)
	tenantHandler := NewTenantHandler(tenantService)
	userHandler := NewUserHandler(userService)

	// Create module handlers
	accountHandler := accounting.NewAccountHandler(accountService)
	journalEntryHandler := accounting.NewJournalEntryHandler(journalEntryService)
	productHandler := inventory.NewProductHandler(productService)
	inventoryTransactionHandler := inventory.NewInventoryTransactionHandler(inventoryTransactionService, productService)
	customerHandler := crm.NewCustomerHandler(customerService, contactService)
	contactHandler := crm.NewContactHandler(contactService, customerService)
	interactionHandler := crm.NewInteractionHandler(interactionService, customerService)

	// Public routes
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/tenants/{subdomain}", tenantHandler.GetTenantBySubdomain).Methods("GET")

	// Admin routes (no tenant context)
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(auth.JWTMiddleware(jwtService))
	adminRouter.HandleFunc("/tenants", tenantHandler.ListTenants).Methods("GET")
	adminRouter.HandleFunc("/tenants", tenantHandler.CreateTenant).Methods("POST")
	adminRouter.HandleFunc("/tenants/{id}", tenantHandler.GetTenant).Methods("GET")
	adminRouter.HandleFunc("/tenants/{id}", tenantHandler.UpdateTenant).Methods("PUT")
	adminRouter.HandleFunc("/tenants/{id}", tenantHandler.DeleteTenant).Methods("DELETE")

	// Tenant routes (with tenant context)
	tenantRouter := r.PathPrefix("/api").Subrouter()
	tenantRouter.Use(auth.JWTMiddleware(jwtService))

	// User routes
	tenantRouter.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	tenantRouter.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	tenantRouter.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	tenantRouter.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	tenantRouter.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Accounting routes
	tenantRouter.HandleFunc("/accounting/accounts", accountHandler.ListAccounts).Methods("GET")
	tenantRouter.HandleFunc("/accounting/accounts", accountHandler.CreateAccount).Methods("POST")
	tenantRouter.HandleFunc("/accounting/accounts/{id}", accountHandler.GetAccount).Methods("GET")
	tenantRouter.HandleFunc("/accounting/accounts/{id}", accountHandler.UpdateAccount).Methods("PUT")
	tenantRouter.HandleFunc("/accounting/accounts/{id}", accountHandler.DeleteAccount).Methods("DELETE")

	tenantRouter.HandleFunc("/accounting/journal-entries", journalEntryHandler.ListJournalEntries).Methods("GET")
	tenantRouter.HandleFunc("/accounting/journal-entries", journalEntryHandler.CreateJournalEntry).Methods("POST")
	tenantRouter.HandleFunc("/accounting/journal-entries/{id}", journalEntryHandler.GetJournalEntry).Methods("GET")
	tenantRouter.HandleFunc("/accounting/journal-entries/{id}", journalEntryHandler.UpdateJournalEntry).Methods("PUT")
	tenantRouter.HandleFunc("/accounting/journal-entries/{id}", journalEntryHandler.DeleteJournalEntry).Methods("DELETE")

	// Inventory routes
	tenantRouter.HandleFunc("/inventory/products", productHandler.ListProducts).Methods("GET")
	tenantRouter.HandleFunc("/inventory/products", productHandler.CreateProduct).Methods("POST")
	tenantRouter.HandleFunc("/inventory/products/{id}", productHandler.GetProduct).Methods("GET")
	tenantRouter.HandleFunc("/inventory/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	tenantRouter.HandleFunc("/inventory/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	tenantRouter.HandleFunc("/inventory/transactions", inventoryTransactionHandler.ListTransactions).Methods("GET")
	tenantRouter.HandleFunc("/inventory/transactions", inventoryTransactionHandler.CreateTransaction).Methods("POST")
	tenantRouter.HandleFunc("/inventory/transactions/{id}", inventoryTransactionHandler.GetTransaction).Methods("GET")
	tenantRouter.HandleFunc("/inventory/transactions/product/{productId}", inventoryTransactionHandler.ListTransactionsByProduct).Methods("GET")

	// CRM routes
	tenantRouter.HandleFunc("/crm/customers", customerHandler.ListCustomers).Methods("GET")
	tenantRouter.HandleFunc("/crm/customers", customerHandler.CreateCustomer).Methods("POST")
	tenantRouter.HandleFunc("/crm/customers/{id}", customerHandler.GetCustomer).Methods("GET")
	tenantRouter.HandleFunc("/crm/customers/{id}", customerHandler.UpdateCustomer).Methods("PUT")
	tenantRouter.HandleFunc("/crm/customers/{id}", customerHandler.DeleteCustomer).Methods("DELETE")

	tenantRouter.HandleFunc("/crm/contacts", contactHandler.CreateContact).Methods("POST")
	tenantRouter.HandleFunc("/crm/contacts/{id}", contactHandler.GetContact).Methods("GET")
	tenantRouter.HandleFunc("/crm/contacts/{id}", contactHandler.UpdateContact).Methods("PUT")
	tenantRouter.HandleFunc("/crm/contacts/{id}", contactHandler.DeleteContact).Methods("DELETE")
	tenantRouter.HandleFunc("/crm/customers/{customerId}/contacts", contactHandler.ListContactsByCustomer).Methods("GET")

	tenantRouter.HandleFunc("/crm/interactions", interactionHandler.CreateInteraction).Methods("POST")
	tenantRouter.HandleFunc("/crm/interactions/{id}", interactionHandler.GetInteraction).Methods("GET")
	tenantRouter.HandleFunc("/crm/interactions/{id}", interactionHandler.UpdateInteraction).Methods("PUT")
	tenantRouter.HandleFunc("/crm/interactions/{id}", interactionHandler.DeleteInteraction).Methods("DELETE")
	tenantRouter.HandleFunc("/crm/customers/{customerId}/interactions", interactionHandler.ListInteractionsByCustomer).Methods("GET")

	// Add CORS middleware
	r.Use(corsMiddleware)

	return router
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}