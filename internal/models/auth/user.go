// Package auth models for auth
package auth

// User - Структура пользователя.
type User struct {
	ID    string `json:"id" binding:"required"`
	Token Tokens `json:"token" binding:"required"`
}
