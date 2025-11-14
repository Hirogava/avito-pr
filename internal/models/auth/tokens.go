// Package auth models for auth
package auth

import "time"

// RefreshToken - Структура токена.
type RefreshToken struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Tokens - Структура токенов.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token" binding:"omitempty"`
}
