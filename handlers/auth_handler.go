package handlers

import (
	"net/http"
	"strings"
	"time"

	"auth-backend/config"
	"auth-backend/models"
	"auth-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	Cfg config.Config
}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func NewAuthHandler(db *gorm.DB, cfg config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	var existing models.User
	if err := h.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		utils.Error(c, http.StatusConflict, "email already registered")
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := models.User{
		FullName:     strings.TrimSpace(req.FullName),
		Email:        email,
		Phone:        strings.TrimSpace(req.Phone),
		PasswordHash: passwordHash,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	token, err := utils.GenerateAccessToken(user.ID, user.Email, h.Cfg.JWTSecret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.Success(c, http.StatusCreated, gin.H{
		"message": "register successful",
		"token":   token,
		"user":    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := utils.GenerateAccessToken(user.ID, user.Email, h.Cfg.JWTSecret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user":    user,
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		// Security: same response for existing/non-existing email.
		utils.Success(c, http.StatusOK, gin.H{"message": "if account exists, reset token generated"})
		return
	}

	token, err := utils.GenerateRandomToken()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate reset token")
		return
	}

	reset := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Used:      false,
	}
	if err := h.DB.Create(&reset).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to save reset token")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{
		"message":     "if account exists, reset token generated",
		"reset_token": token,
		"note":        "in production this token should be sent by email",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var reset models.PasswordResetToken
	err := h.DB.Where("token = ? AND used = ?", strings.TrimSpace(req.Token), false).First(&reset).Error
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid reset token")
		return
	}
	if time.Now().After(reset.ExpiresAt) {
		utils.Error(c, http.StatusBadRequest, "reset token expired")
		return
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	if err := h.DB.Model(&models.User{}).Where("id = ?", reset.UserID).Update("password_hash", passwordHash).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to update password")
		return
	}
	if err := h.DB.Model(&reset).Update("used", true).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to close reset token")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "password reset successful"})
}
