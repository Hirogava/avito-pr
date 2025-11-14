package postgres

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/models/types"
)

func TestCreatePullRequestSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.PullRequestCreateRequest{
		PullRequestID:   "pr-1",
		PullRequestName: "Add feature",
		AuthorID:        "author-1",
	}

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(req.PullRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectQuery(`SELECT team_name FROM users WHERE user_id = \$1 AND is_active = TRUE`).
		WithArgs(req.AuthorID).
		WillReturnRows(sqlmock.NewRows([]string{"team_name"}).AddRow("backend"))

	mock.ExpectQuery(`SELECT user_id FROM users\s+WHERE team_name = \$1`).
		WithArgs("backend", req.AuthorID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("reviewer-1"))

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO pull_requests`).
		WithArgs(req.PullRequestID, req.PullRequestName, req.AuthorID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO pr_reviewers`).
		WithArgs(req.PullRequestID, "reviewer-1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	pr, err := manager.CreatePullRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.PullRequestID != req.PullRequestID {
		t.Fatalf("unexpected response %#v", pr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCreatePullRequestExists(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.PullRequestCreateRequest{PullRequestID: "pr-1"}

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(req.PullRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	_, err := manager.CreatePullRequest(req)
	if !errors.Is(err, dbErrors.ErrorPRAlreadyExists) {
		t.Fatalf("expected ErrorPRAlreadyExists, got %v", err)
	}
}

func TestMergePullRequestSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.PullRequestMergeRequest{PullRequestID: "pr-1"}
	mock.ExpectQuery(`SELECT pull_request_id`).
		WithArgs(req.PullRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
			AddRow(req.PullRequestID, "Feature", "author", "OPEN"))

	mock.ExpectExec(`UPDATE pull_requests SET status = 'MERGED', merged_at = NOW\(\)`).
		WithArgs(req.PullRequestID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`SELECT reviewer_id FROM pr_reviewers`).
		WithArgs(req.PullRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"reviewer_id"}).AddRow("rev-1"))

	resp, err := manager.MergePullRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != types.PRStatusMerged {
		t.Fatalf("expected status %s, got %s", types.PRStatusMerged, resp.Status)
	}
	if len(resp.AssignedReviewers) != 1 || resp.AssignedReviewers[0] != "rev-1" {
		t.Fatalf("unexpected reviewers %+v", resp.AssignedReviewers)
	}
}

func TestMergePullRequestNotFound(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.PullRequestMergeRequest{PullRequestID: "missing"}

	mock.ExpectQuery(`SELECT pull_request_id`).
		WithArgs(req.PullRequestID).
		WillReturnError(sql.ErrNoRows)

	_, err := manager.MergePullRequest(req)
	if !errors.Is(err, dbErrors.ErrorPRSNotFound) {
		t.Fatalf("expected ErrorPRSNotFound, got %v", err)
	}
}

func TestReassignPRAuthorSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.PullRequestReassignRequest{
		PullRequestID: "pr-1",
		OldUserID:     "old",
	}

	mock.ExpectQuery(`SELECT status, author_id FROM pull_requests`).
		WithArgs(req.PullRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"status", "author_id"}).AddRow("OPEN", "author"))

	mock.ExpectQuery(`SELECT EXISTS\(`).
		WithArgs(req.PullRequestID, req.OldUserID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(`SELECT team_name FROM users WHERE user_id = \$1`).
		WithArgs("author").
		WillReturnRows(sqlmock.NewRows([]string{"team_name"}).AddRow("backend"))

	mock.ExpectQuery(`SELECT user_id FROM users`).
		WithArgs("backend", req.OldUserID, "author").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("new-reviewer"))

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM pr_reviewers`).
		WithArgs(req.PullRequestID, req.OldUserID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO pr_reviewers`).
		WithArgs(req.PullRequestID, "new-reviewer").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(`SELECT reviewer_id FROM pr_reviewers`).
		WithArgs(req.PullRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"reviewer_id"}).AddRow("new-reviewer"))

	resp, err := manager.ReassignPRAuthor(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ReplacedBy != "new-reviewer" {
		t.Fatalf("expected replacement new-reviewer, got %s", resp.ReplacedBy)
	}
}

func TestReassignPRAuthorNoCandidate(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.PullRequestReassignRequest{
		PullRequestID: "pr-1",
		OldUserID:     "old",
	}

	mock.ExpectQuery(`SELECT status, author_id FROM pull_requests`).
		WithArgs(req.PullRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"status", "author_id"}).AddRow("OPEN", "author"))

	mock.ExpectQuery(`SELECT EXISTS\(`).
		WithArgs(req.PullRequestID, req.OldUserID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(`SELECT team_name FROM users WHERE user_id = \$1`).
		WithArgs("author").
		WillReturnRows(sqlmock.NewRows([]string{"team_name"}).AddRow("backend"))

	mock.ExpectQuery(`SELECT user_id FROM users`).
		WithArgs("backend", req.OldUserID, "author").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))

	_, err := manager.ReassignPRAuthor(req)
	if !errors.Is(err, dbErrors.ErrorNoCandidateForReviewer) {
		t.Fatalf("expected ErrorNoCandidateForReviewer, got %v", err)
	}
}
