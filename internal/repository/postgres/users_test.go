package postgres

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
)

func TestSetUserIsActiveSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.UserSetIsActiveRequest{
		UserID:   "user-1",
		IsActive: true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET is_active = $1 WHERE user_id = $2 RETURNING is_active, username, team_name, user_id`)).
		WithArgs(req.IsActive, req.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"is_active", "username", "team_name", "user_id"}).
			AddRow(true, "alice", "backend", req.UserID))

	user, err := manager.SetUserIsActive(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.UserID != req.UserID || !user.IsActive {
		t.Fatalf("unexpected user response: %#v", user)
	}
}

func TestSetUserIsActiveNotFound(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.UserSetIsActiveRequest{UserID: "missing"}

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET is_active = $1 WHERE user_id = $2 RETURNING is_active, username, team_name, user_id`)).
		WithArgs(req.IsActive, req.UserID).
		WillReturnError(sql.ErrNoRows)

	_, err := manager.SetUserIsActive(req)
	if err != dbErrors.ErrorUserNotFound {
		t.Fatalf("expected ErrorUserNotFound, got %v", err)
	}
}

func TestGetUsersReviewSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.UsersGetReviewQuery{UserID: "user"}

	rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
		AddRow("pr1", "Fix bug", "author", "OPEN")

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			pr.author_id,
			pr.status
		FROM pr_reviewers r
		JOIN pull_requests pr ON r.pull_request_id = pr.pull_request_id
		WHERE r.reviewer_id = $1
	`)).
		WithArgs(req.UserID).
		WillReturnRows(rows)

	review, err := manager.GetUsersReview(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(review.PullRequests) != 1 {
		t.Fatalf("expected 1 PR, got %d", len(review.PullRequests))
	}
}

func TestGetUsersReviewQueryError(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.UsersGetReviewQuery{UserID: "user"}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			pr.author_id,
			pr.status
		FROM pr_reviewers r
		JOIN pull_requests pr ON r.pull_request_id = pr.pull_request_id
		WHERE r.reviewer_id = $1
	`)).
		WithArgs(req.UserID).
		WillReturnError(sql.ErrConnDone)

	_, err := manager.GetUsersReview(req)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetUsersSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"username", "team_name", "user_id", "is_active"}).
		AddRow("alice", "backend", "u1", true).
		AddRow("bob", "backend", "u2", false)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT username, team_name, user_id, is_active FROM users
		`)).
		WillReturnRows(rows)

	users, err := manager.GetUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestGetUsersScanError(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"username", "team_name", "user_id", "is_active"}).
		AddRow(nil, "backend", "u1", true)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT username, team_name, user_id, is_active FROM users
		`)).
		WillReturnRows(rows)

	_, err := manager.GetUsers()
	if err == nil {
		t.Fatal("expected scan error")
	}
}
