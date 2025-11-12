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

var secret = os.Getenv("JWT_SECRET")

func ParseToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, jwt.ErrTokenMalformed
	} else {
		return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	}
}

func GenerateRefreshToken(manager *postgres.Manager, userId string) (string, error) {
	logger.Logger.Debug("Generating refresh token", "user_id", userId)

	token := uuid.New().String()

	err := manager.SaveRefreshToken(token, userId)
	if err != nil {
		logger.Logger.Error("Failed to save refresh token", "user_id", userId, "error", err.Error())
		return "", err
	}

	logger.Logger.Debug("Refresh token generated and saved", "user_id", userId)
	return token, nil
}

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

func AddAccessTime() int64 {
	return time.Now().Add(15 * time.Minute).Unix()
}

func ValidateRefreshToken(manager *postgres.Manager, userId string) (authModels.RefreshToken, error) {
	token, err := manager.GetRefreshToken(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return authModels.RefreshToken{}, jwt.ErrTokenExpired
		}
		return token, err
	}

	if time.Now().After(token.ExpiredAt) {
		if err := manager.DeleteRefreshToken(userId, token.Token); err != nil {
			return authModels.RefreshToken{}, err
		} else {
			return authModels.RefreshToken{}, jwt.ErrTokenExpired
		}
	}
	return token, nil
}

func GetClaims(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return claims, err
}
