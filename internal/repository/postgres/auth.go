package postgres

import (
	"time"

	authModels "github.com/Hirogava/avito-pr/internal/models/auth"
)

func (manager *Manager) SaveRefreshToken(token string, userId string) error {
	expiredAt := time.Now().Add(time.Hour * 24 * 7)

	_, err := manager.Conn.Exec(`INSERT INTO session (user_id, token, expires_at) VALUES ($1, $2, $3)`, userId, token, expiredAt)
	return err
}

func (manager *Manager) GetRefreshToken(userId string) (authModels.RefreshToken, string, error) {
	var token authModels.RefreshToken

	if err := manager.Conn.QueryRow(`SELECT id, token, expires_at FROM session WHERE user_id = $1`, userId).Scan(&token.ID, &token.Token, &token.ExpiredAt); err != nil {
		return authModels.RefreshToken{}, "", err
	}

	role, err := manager.GetUserRoleByID(userId)
	if err != nil {
		return authModels.RefreshToken{}, "", err
	}

	return token, role, nil
}

func (manager *Manager) GetUserRoleByID(userId string) (string, error) {
	var role string
	var isAdmin bool

	if err := manager.Conn.QueryRow(`SELECT is_admin FROM users WHERE id = $1`, userId).Scan(&isAdmin); err != nil {
		return "", err
	}

	if isAdmin {
		role = "admin"
	} else {
		role = "user"
	}

	return role, nil
}

func (manager *Manager) DeleteRefreshToken(userId string, token string) error {
	_, err := manager.Conn.Exec(`DELETE FROM session WHERE user_id = $1 and token = $2`, userId, token)

	return err
}

func (manager *Manager) UpdateRefreshToken(userId string, token string) error {
	if _, err := manager.Conn.Exec(`UPDATE session SET expires_at = $1, token = $2 WHERE user_id = $3`, time.Now().Add(time.Hour*24*7), token, userId); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) SetUserStatus(userId string, isAdmin bool) error {
	_, err := manager.Conn.Exec(`UPDATE users SET is_admin = $1 WHERE id = $2`, isAdmin, userId)
	return err
}
