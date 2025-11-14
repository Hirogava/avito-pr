// Package postgres implements the repository interface for PostgreSQL.
package postgres

import (
	"database/sql"
	"time"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
)

// CreateTeam - создает новую команду
func (m *Manager) CreateTeam(req reqres.TeamAddRequest) (*reqres.TeamResponse, error) {
	tx, err := m.Conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	var exists bool
	err = tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM teams WHERE team_name = $1)`, req.TeamName).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, dbErrors.ErrorTeamAlreadyExists
	}

	_, err = tx.Exec(`INSERT INTO teams (team_name, created_at) VALUES ($1, $2)`, req.TeamName, time.Now())
	if err != nil {
		return nil, err
	}

	for _, member := range req.Members {
		_, err := tx.Exec(`
			INSERT INTO users (user_id, username, team_name, is_active, created_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (user_id) DO UPDATE
			SET username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active,
				updated_at = NOW();
		`, member.UserID, member.Username, req.TeamName, member.IsActive)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &reqres.TeamResponse{
		TeamName: req.TeamName,
		Members:  req.Members,
	}, nil
}

// GetTeam - возвращает команду
func (m *Manager) GetTeam(teamName string) (*reqres.TeamResponse, error) {
	rows, err := m.Conn.Query(`
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username;
	`, teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dbErrors.ErrorTeamNotFound
		}

		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var members []reqres.TeamMemberResponse
	for rows.Next() {
		var mbr reqres.TeamMemberResponse
		if err := rows.Scan(&mbr.UserID, &mbr.Username, &mbr.IsActive); err != nil {
			return nil, err
		}
		members = append(members, mbr)
	}

	return &reqres.TeamResponse{
		TeamName: teamName,
		Members:  members,
	}, nil
}
