package middleware

import (
	"net/http"
	"strings"

	"auth-backend/config"
	"auth-backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthRequired(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, http.StatusUnauthorized, "missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Error(c, http.StatusUnauthorized, "invalid Authorization format")
			c.Abort()
			return
		}

		claims, err := utils.ParseAccessToken(parts[1], cfg.JWTSecret)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
