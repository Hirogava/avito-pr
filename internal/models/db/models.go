// Package db models for database
package db

import (
	"database/sql"
	"time"

	"github.com/Hirogava/avito-pr/internal/models/types"
)

// UserDBModel - Модель пользователя в базе данных.
type UserDBModel struct {
	UserID    string       `db:"user_id"`
	Username  string       `db:"username"`
	TeamName  string       `db:"team_name"`
	IsActive  bool         `db:"is_active"`
	IsAdmin   bool         `db:"is_admin"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// TeamDBModel - Модель команды в базе данных.
type TeamDBModel struct {
	TeamName  string       `db:"team_name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// PullRequestDBModel - Модель PR в базе данных.
type PullRequestDBModel struct {
	PullRequestID   string         `db:"pull_request_id"`
	PullRequestName string         `db:"pull_request_name"`
	AuthorID        string         `db:"author_id"`
	Status          types.PRStatus `db:"status"`
	CreatedAt       time.Time      `db:"created_at"`
	MergedAt        sql.NullTime   `db:"merged_at"`
}

// PullRequestReviewerDBModel - Модель участника PR в базе данных.
type PullRequestReviewerDBModel struct {
	PullRequestID string    `db:"pull_request_id"`
	ReviewerID    string    `db:"reviewer_id"`
	AssignedAt    time.Time `db:"assigned_at"`
}
