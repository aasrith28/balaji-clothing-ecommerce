package main

import (
	"log"
	"os"

	"github.com/balaji/ecommerce/internal/config"
	"github.com/balaji/ecommerce/internal/database"
	"github.com/balaji/ecommerce/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Connect DB
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Migrate
	if err := database.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Seed
	if err := database.Seed(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	// Gin setup
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.Setup(r, cfg)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Balaji Shop server running on http://localhost:%s", port)
	log.Printf("📧 Admin: admin@balaji.com / admin123")
	log.Printf("👤 User:  user@balaji.com / user123")

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
