package middleware

import (
	"auth-backend/config"
	"auth-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := utils.ParseAccessToken(parts[1], cfg.AccessTokenSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired access token"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
