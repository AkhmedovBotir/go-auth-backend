package main

import (
	"log"

	"auth-backend/config"
	"auth-backend/database"
	"auth-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabasePath)

	router := gin.Default()
	routes.Register(router, db, cfg)

	log.Printf("Server running on %s", cfg.Port)
	if err := router.Run(cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
