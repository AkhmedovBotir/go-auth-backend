package database

import (
	"log"

	"auth-backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Connect(path string) *gorm.DB {
	// Using pure-Go sqlite driver to avoid CGO_ENABLED requirements.
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.PasswordResetToken{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
