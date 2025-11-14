// Package auth provides functions for working with JWT tokens
package auth

import (
	"database/sql"
	"os"
	"time"

	"github.com/Hirogava/avito-pr/internal/config/logger"
	authModels "github.com/Hirogava/avito-pr/internal/models/auth"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Секретный ключ для JWT
var secret = os.Getenv("JWT_SECRET")

// ParseToken проверяет токен на валидность
func ParseToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, jwt.ErrTokenMalformed
	}

	return jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

// GenerateRefreshToken генерирует токен для пользователя
func GenerateRefreshToken(manager *postgres.Manager, userID string) (string, error) {
	logger.Logger.Debug("Generating refresh token", "user_id", userID)

	token := uuid.New().String()

	err := manager.SaveRefreshToken(token, userID)
	if err != nil {
		logger.Logger.Error("Failed to save refresh token", "user_id", userID, "error", err.Error())
		return "", err
	}

	logger.Logger.Debug("Refresh token generated and saved", "user_id", userID)
	return token, nil
}

// GenerateAccessToken генерирует токен доступа для пользователя
func GenerateAccessToken(claims jwt.MapClaims) (string, error) {
	logger.Logger.Debug("Generating access token", "user_id", claims["id"])

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Logger.Error("Failed to sign access token", "user_id", claims["id"], "error", err.Error())
		return "", err
	}

	logger.Logger.Debug("Access token generated successfully", "user_id", claims["id"])
	return accessToken, nil
}

// AddAccessTime добавляет 15 минут к текущему времени
func AddAccessTime() int64 {
	return time.Now().Add(15 * time.Minute).Unix()
}

// ValidateRefreshToken проверяет токен на валидность
func ValidateRefreshToken(manager *postgres.Manager, userID string) (authModels.RefreshToken, string, error) {
	token, role, err := manager.GetRefreshToken(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return authModels.RefreshToken{}, "", jwt.ErrTokenExpired
		}
		return token, role, err
	}

	if time.Now().After(token.ExpiredAt) {
		if err := manager.DeleteRefreshToken(userID, token.Token); err != nil {
			return authModels.RefreshToken{}, "", err
		}

		return authModels.RefreshToken{}, "", jwt.ErrTokenExpired
	}
	return token, role, nil
}

// GetClaims извлекает данные из токена
func GetClaims(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return claims, err
}
