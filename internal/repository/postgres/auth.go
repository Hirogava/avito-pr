// Package postgres implements the repository interface for PostgreSQL.
package postgres

import (
	"time"

	authModels "github.com/Hirogava/avito-pr/internal/models/auth"
)

// SaveRefreshToken - сохраняет refresh token в базе данных.
func (manager *Manager) SaveRefreshToken(token string, userID string) error {
	expiredAt := time.Now().Add(time.Hour * 24 * 7)

	_, err := manager.Conn.Exec(`INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)`, userID, token, expiredAt)
	return err
}

// GetRefreshToken - получает refresh token из базы данных.
func (manager *Manager) GetRefreshToken(userID string) (authModels.RefreshToken, string, error) {
	var token authModels.RefreshToken

	if err := manager.Conn.QueryRow(`SELECT id, token, expires_at FROM sessions WHERE user_id = $1`, userID).Scan(&token.ID, &token.Token, &token.ExpiredAt); err != nil {
		return authModels.RefreshToken{}, "", err
	}

	role, err := manager.GetUserRoleByID(userID)
	if err != nil {
		return authModels.RefreshToken{}, "", err
	}

	return token, role, nil
}

// GetUserRoleByID - получает роль пользователя по его ID.
func (manager *Manager) GetUserRoleByID(userID string) (string, error) {
	var role string
	var isAdmin bool

	if err := manager.Conn.QueryRow(`SELECT is_admin FROM users WHERE user_id = $1`, userID).Scan(&isAdmin); err != nil {
		return "", err
	}

	if isAdmin {
		role = "admin"
	} else {
		role = "user"
	}

	return role, nil
}

// DeleteRefreshToken - удаляет refresh token из базы данных.
func (manager *Manager) DeleteRefreshToken(userID string, token string) error {
	_, err := manager.Conn.Exec(`DELETE FROM sessions WHERE user_id = $1 and token = $2`, userID, token)

	return err
}

// UpdateRefreshToken - обновляет refresh token в базе данных.
func (manager *Manager) UpdateRefreshToken(userID string, token string) error {
	if _, err := manager.Conn.Exec(`UPDATE sessions SET expires_at = $1, token = $2 WHERE user_id = $3`, time.Now().Add(time.Hour*24*7), token, userID); err != nil {
		return err
	}

	return nil
}

// SetUserStatus - устанавливает статус пользователя в базе данных.
func (manager *Manager) SetUserStatus(userID string, isAdmin bool) error {
	_, err := manager.Conn.Exec(`UPDATE users SET is_admin = $1 WHERE user_id = $2`, isAdmin, userID)
	return err
}
