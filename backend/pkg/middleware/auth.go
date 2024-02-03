package middleware

import (
	"aitu-funpage/backend/pkg/logger"
	"aitu-funpage/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequireAuthMiddleware(c *gin.Context) {
	logger := logger.GetLogger()

	var token string
	authHeader := c.GetHeader("Authorization")
	cookie, err := c.Cookie("jwt")
	if err != nil {
		logger.Error("JWT token not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token not found"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	} else if err == nil {
		token = cookie
	}

	if authHeader == "" {
		logger.Error("Authorization header not found")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Try to sign in first"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if token == "" {
		logger.Error("Invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id, username, userType, err := utils.VerifyToken(token)
	if err != nil {
		logger.Error("Token verification failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token verification failed"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("id", id)
	c.Set("username", username)
	c.Set("userType", userType)

	logger.Info("Authentication successful")
	c.Next()
}
