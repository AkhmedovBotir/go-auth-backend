package handlers

import (
	"net/http"
	"strings"

	"auth-backend/config"
	"auth-backend/models"
	"auth-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	DB  *gorm.DB
	Cfg config.Config
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required,min=2"`
	Phone    string `json:"phone"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

func NewProfileHandler(db *gorm.DB, cfg config.Config) *ProfileHandler {
	return &ProfileHandler{DB: db, Cfg: cfg}
}

func getUserIDFromClaims(c *gin.Context) (uint, bool) {
	claimsValue, exists := c.Get("claims")
	if !exists {
		return 0, false
	}
	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, false
	}
	return uint(sub), true
}

func (h *ProfileHandler) GetMe(c *gin.Context) {
	userID, ok := getUserIDFromClaims(c)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "user not found")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"user": user})
}

func (h *ProfileHandler) UpdateMe(c *gin.Context) {
	userID, ok := getUserIDFromClaims(c)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]any{
		"full_name": strings.TrimSpace(req.FullName),
		"phone":     strings.TrimSpace(req.Phone),
	}
	if err := h.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to update profile")
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "user not found")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{
		"message": "profile updated",
		"user":    user,
	})
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID, ok := getUserIDFromClaims(c)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "user not found")
		return
	}

	if !utils.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		utils.Error(c, http.StatusBadRequest, "current password is incorrect")
		return
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to hash new password")
		return
	}

	if err := h.DB.Model(&user).Update("password_hash", newHash).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to change password")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "password changed successfully"})
}
