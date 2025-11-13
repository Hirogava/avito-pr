package reqres

// TeamAddRequest - Запрос на создание/обновление команды.
type TeamAddRequest struct {
	TeamName string               `json:"team_name" binding:"required"`
	Members  []TeamMemberResponse `json:"members" binding:"required,min=1"`
}

// UserSetIsActiveRequest - Запрос на установку флага активности пользователя.
type UserSetIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// PullRequestCreateRequest - Запрос на создание PR.
type PullRequestCreateRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

// PullRequestMergeRequest - Запрос на мерж PR.
type PullRequestMergeRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

// PullRequestReassignRequest - Запрос на переназначение ревьювера.
type PullRequestReassignRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_reviewer_id" binding:"required"`
}

// SetAdminRequest - Запрос на установку флага админа пользователя.
type SetAdminRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	IsAdmin bool   `json:"is_admin"`
}
