package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yookibooki/erp/internal/api"
	"github.com/yookibooki/erp/internal/auth"
	"github.com/yookibooki/erp/internal/config"
	"github.com/yookibooki/erp/internal/db"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close()

	// Create repositories
	tenantRepo := db.NewTenantRepository(database)
	userRepo := db.NewUserRepository(database)
	
	// Create module repositories
	accountRepo := db.NewAccountRepository(database)
	journalEntryRepo := db.NewJournalEntryRepository(database)
	productRepo := db.NewProductRepository(database)
	inventoryTransactionRepo := db.NewInventoryTransactionRepository(database)
	customerRepo := db.NewCustomerRepository(database)
	contactRepo := db.NewContactRepository(database)
	interactionRepo := db.NewInteractionRepository(database)

	// Create JWT service
	jwtService := auth.NewJWTService(cfg.JWT)

	// Create router
	router := api.NewRouter(
		tenantRepo,
		userRepo,
		accountRepo,
		journalEntryRepo,
		productRepo,
		inventoryTransactionRepo,
		customerRepo,
		contactRepo,
		interactionRepo,
		jwtService,
	)

	// Start server
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}