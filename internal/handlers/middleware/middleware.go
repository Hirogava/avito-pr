package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/Hirogava/avito-pr/internal/config/logger"
	"github.com/Hirogava/avito-pr/internal/service/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		logger.Logger.Debug("Auth middleware processing request",
			"method", method,
			"path", path,
			"ip", c.ClientIP())

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			logger.Logger.Warn("Missing or invalid Authorization header",
				"method", method,
				"path", path,
				"ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth token required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := auth.ParseToken(tokenString)
		if err != nil {
			logger.Logger.Warn("Failed to parse JWT token",
				"method", method,
				"path", path,
				"ip", c.ClientIP(),
				"error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth token required"})
			c.Abort()
			return
		}
		if !token.Valid {
			logger.Logger.Warn("Invalid JWT token",
				"method", method,
				"path", path,
				"ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				logger.Logger.Warn("JWT token expired",
					"method", method,
					"path", path,
					"ip", c.ClientIP(),
					"exp", int64(exp))
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
		}

		idString, ok := claims["id"].(string)
		if !ok {
			logger.Logger.Warn("Invalid token claims - missing user ID",
				"method", method,
				"path", path,
				"ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || (role != "user" && role != "admin") {
			logger.Logger.Warn("Insufficient permissions",
				"method", method,
				"path", path,
				"user_id", idString,
				"ip", c.ClientIP())
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		logger.Logger.Debug("Authentication successful",
			"method", method,
			"path", path,
			"user_id", idString,
			"ip", c.ClientIP())

		c.Set("role", role)
		c.Set("userID", idString)
		c.Next()
	}
}
