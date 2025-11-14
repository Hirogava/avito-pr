// Package auth provides handlers for auth
package auth

import (
	"net/http"

	"github.com/Hirogava/avito-pr/internal/config/logger"
	"github.com/Hirogava/avito-pr/internal/handlers/middleware"
	authModels "github.com/Hirogava/avito-pr/internal/models/auth"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"
	tokens "github.com/Hirogava/avito-pr/internal/service/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// InitAuthHandlers - инициализация роутов для авторизации
func InitAuthHandlers(r *gin.Engine, manager *postgres.Manager) {
	v1 := r.Group("/auth")
	{
		v1.POST("/admin", func(c *gin.Context) {
			Admin(c, manager)
		})
	}

	secureV1 := r.Group("/auth")
	secureV1.Use(middleware.AuthMiddleware())
	{
		secureV1.POST("/refresh", func(c *gin.Context) {
			RefreshToken(c, manager)
		})
	}
}

// Admin - авторизация роли пользователя
func Admin(c *gin.Context, manager *postgres.Manager) {
	logger.Logger.Info("Login attempt", "ip", c.ClientIP())

	var req reqres.SetAdminRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid login request", "ip", c.ClientIP(), "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Debug("Processing login", "user_id", req.UserID, "ip", c.ClientIP())

	err := manager.SetUserStatus(req.UserID, req.IsAdmin)
	if err != nil {
		logger.Logger.Error("Failed to set user admin status", "ip", c.ClientIP(), "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var role string
	c.Set("userID", req.UserID)
	if req.IsAdmin {
		c.Set("role", "admin")
		role = "admin"
	} else {
		c.Set("role", "user")
		role = "user"
	}

	var refreshToken string
	token, _, err := tokens.ValidateRefreshToken(manager, req.UserID)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			logger.Logger.Debug("Refresh token expired, generating new one", "user_id", req.UserID)
			refreshToken, err = tokens.GenerateRefreshToken(manager, req.UserID)
			if err != nil {
				logger.Logger.Error("Failed to generate refresh token", "user_id", req.UserID, "error", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			logger.Logger.Error("Failed to validate refresh token", "user_id", req.UserID, "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if token.Token == "" {
		token.Token = refreshToken
	}

	claims := jwt.MapClaims{
		"id":   req.UserID,
		"exp":  tokens.AddAccessTime(),
		"role": role,
	}

	var accessToken string

	if accessToken, err = tokens.GenerateAccessToken(claims); err != nil {
		logger.Logger.Error("Failed to generate access token", "user_id", req.UserID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}

	logger.Logger.Info("Login successful", "user_id", req.UserID, "ip", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"userID":        req.UserID,
		"role":          role,
		"access_token":  accessToken,
		"refresh_token": token.Token,
	})
}

// RefreshToken - обновление токена
func RefreshToken(c *gin.Context, manager *postgres.Manager) {
	userID := c.GetString("userID")
	logger.Logger.Info("Token refresh attempt", "user_id", userID, "ip", c.ClientIP())

	var t authModels.Tokens
	if err := c.ShouldBindJSON(&t); err != nil {
		logger.Logger.Warn("Invalid refresh token request", "user_id", userID, "ip", c.ClientIP(), "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, role, err := tokens.ValidateRefreshToken(manager, userID)
	if err != nil {
		switch err {
		case jwt.ErrTokenExpired:
			logger.Logger.Warn("Refresh token expired", "user_id", userID, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
			return
		default:
			logger.Logger.Error("Failed to validate refresh token", "user_id", userID, "ip", c.ClientIP(), "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	t.RefreshToken = refreshToken.Token

	claims := jwt.MapClaims{
		"id":   userID,
		"exp":  tokens.AddAccessTime(),
		"role": role,
	}

	if t.AccessToken, err = tokens.GenerateAccessToken(claims); err != nil {
		logger.Logger.Error("Failed to generate new access token", "user_id", userID, "ip", c.ClientIP(), "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Info("Token refresh successful", "user_id", userID, "ip", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"access_token":  t.AccessToken,
		"refresh_token": t.RefreshToken,
	})
}
