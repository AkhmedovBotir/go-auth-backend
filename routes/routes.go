package routes

import (
	"auth-backend/config"
	"auth-backend/handlers"
	"auth-backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(router *gin.Engine, db *gorm.DB, cfg config.Config) {
	authHandler := handlers.NewAuthHandler(db, cfg)
	profileHandler := handlers.NewProfileHandler(db, cfg)

	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	profile := api.Group("/profile")
	profile.Use(middleware.AuthRequired(cfg))
	{
		profile.GET("/me", profileHandler.GetMe)
		profile.PUT("/me", profileHandler.UpdateMe)
		profile.PUT("/change-password", profileHandler.ChangePassword)
	}
}
