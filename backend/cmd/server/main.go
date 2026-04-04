package main

import (
	"log"

	"github.com/wenaixi/wenxi-cloud/backend/config"
	"github.com/wenaixi/wenxi-cloud/backend/internal/api"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/database"
	"github.com/wenaixi/wenxi-cloud/backend/internal/repository"
)

func main() {
	cfg := config.Load()

	// Initialize database (GORM with auto-migration)
	db, err := database.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Printf("Database initialized: %s", cfg.DBPath)

	// Inject DB into router's global store
	api.SetDB(db)

	// Initialize repositories (ensures DB is connected before first request)
	_ = repository.NewUserRepository(db)
	_ = repository.NewFileRepository(db)
	_ = repository.NewFolderRepository(db)
	_ = repository.NewShareRepository(db)
	_ = repository.NewLanZouTokenRepository(db)
	_ = repository.NewUploadSessionRepository(db)
	_ = repository.NewRecycleBinRepository(db)

	// Setup router (modular architecture)
	r := api.SetupRouter(cfg)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
