package reqres

import (
	"time"

	"github.com/Hirogava/avito-pr/internal/models/types"
)

// TeamMemberResponse - Модель участника команды для ответа API.
type TeamMemberResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// TeamResponse - Модель команды для ответа API.
type TeamResponse struct {
	TeamName string               `json:"team_name"`
	Members  []TeamMemberResponse `json:"members"`
}

// UserResponse - Модель пользователя для ответа API.
type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// PullRequestResponse - Полная модель PR для ответа API.
type PullRequestResponse struct {
	PullRequestID    string    		`json:"pull_request_id"`
	PullRequestName  string    		`json:"pull_request_name"`
	AuthorID         string    		`json:"author_id"`
	Status           types.PRStatus `db:"status"`
	AssignedReviewers []string  	`json:"assigned_reviewers"`
	CreatedAt        time.Time 		`json:"createdAt,omitempty"`
	MergedAt         *time.Time 	`json:"mergedAt,omitempty"`
}

// PullRequestShortResponse - Укороченная модель PR для ответа API.
type PullRequestShortResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// PullRequestMiddleResponse - Средняя модель PR для ответа API.
type PullRequestMiddleResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
	AssignedReviewers []string  	`json:"assigned_reviewers"`
}

// PullRequestListResponse - Модель списка PR для ответа API.
type PullRequestListResponse struct {
	UserID	  string                     `json:"user_id"`
	PullRequests []PullRequestShortResponse `json:"pull_requests"`
}

// PullRequestReassignResponse - Модель ответа на переназначение ревьювера.
type PullRequestReassignResponse struct {
	PR         PullRequestMiddleResponse `json:"pull_request"`
	ReplacedBy string                  `json:"replaced_by"`
}

// ErrorResponse - Модель ошибки для ответа API.
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
