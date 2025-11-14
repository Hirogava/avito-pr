// Package postgres implements the repository interface for PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/models/types"
)

// CreatePullRequest - создает PR с двумя случайными ревьюверами
func (m *Manager) CreatePullRequest(req reqres.PullRequestCreateRequest) (reqres.PullRequestResponse, error) {
	ctx := context.Background()

	var exists bool
	err := m.Conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`, req.PullRequestID).Scan(&exists)
	if err != nil {
		return reqres.PullRequestResponse{}, err
	}
	if exists {
		return reqres.PullRequestResponse{}, dbErrors.ErrorPRAlreadyExists
	}

	var teamName string
	err = m.Conn.QueryRowContext(ctx, `
		SELECT team_name FROM users WHERE user_id = $1 AND is_active = TRUE
	`, req.AuthorID).Scan(&teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reqres.PullRequestResponse{}, dbErrors.ErrorUserNotFound
		}
		return reqres.PullRequestResponse{}, err
	}

	rows, err := m.Conn.QueryContext(ctx, `
		SELECT user_id FROM users 
		WHERE team_name = $1 AND is_active = TRUE AND user_id != $2
	`, teamName, req.AuthorID)
	if err != nil {
		return reqres.PullRequestResponse{}, err
	}
	defer rows.Close() //nolint:errcheck

	var candidates []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return reqres.PullRequestResponse{}, err
		}
		candidates = append(candidates, uid)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	reviewers := []string{}
	if len(candidates) >= 2 {
		reviewers = candidates[:2]
	} else {
		reviewers = candidates
	}

	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		return reqres.PullRequestResponse{}, err
	}
	defer tx.Rollback() //nolint:errcheck

	_, err = tx.ExecContext(ctx, `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
		VALUES ($1, $2, $3, 'OPEN')
	`, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		return reqres.PullRequestResponse{}, err
	}

	for _, rid := range reviewers {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
			VALUES ($1, $2)
		`, req.PullRequestID, rid)
		if err != nil {
			return reqres.PullRequestResponse{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return reqres.PullRequestResponse{}, err
	}

	return reqres.PullRequestResponse{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            types.PRStatusOpen,
		AssignedReviewers: reviewers,
	}, nil
}

// MergePullRequest - мержит PR
func (m *Manager) MergePullRequest(req reqres.PullRequestMergeRequest) (reqres.PullRequestResponse, error) {
	ctx := context.Background()

	var pr reqres.PullRequestResponse
	err := m.Conn.QueryRowContext(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status
		FROM pull_requests WHERE pull_request_id = $1
	`, req.PullRequestID).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reqres.PullRequestResponse{}, dbErrors.ErrorPRSNotFound
		}
		return reqres.PullRequestResponse{}, err
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	_, err = m.Conn.ExecContext(ctx, `
		UPDATE pull_requests SET status = 'MERGED', merged_at = NOW()
		WHERE pull_request_id = $1
	`, req.PullRequestID)
	if err != nil {
		return reqres.PullRequestResponse{}, err
	}

	rows, err := m.Conn.QueryContext(ctx, `
		SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1
	`, req.PullRequestID)
	if err != nil {
		return reqres.PullRequestResponse{}, err
	}
	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return reqres.PullRequestResponse{}, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, rid)
	}

	pr.Status = types.PRStatusMerged
	pr.MergedAt = &time.Time{} //&time.Now() позже поправлю

	return pr, nil
}

// ReassignPRAuthor - меняет автора PR на нового случайного ревьюера
func (m *Manager) ReassignPRAuthor(req reqres.PullRequestReassignRequest) (reqres.PullRequestReassignResponse, error) {
	ctx := context.Background()

	var status, authorID string
	err := m.Conn.QueryRowContext(ctx, `
		SELECT status, author_id FROM pull_requests WHERE pull_request_id = $1
	`, req.PullRequestID).Scan(&status, &authorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reqres.PullRequestReassignResponse{}, dbErrors.ErrorPRSNotFound
		}
		return reqres.PullRequestReassignResponse{}, err
	}

	if status == "MERGED" {
		return reqres.PullRequestReassignResponse{}, dbErrors.ErrorPRMerged
	}

	var assigned bool
	err = m.Conn.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM pr_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2
		)
	`, req.PullRequestID, req.OldUserID).Scan(&assigned)
	if err != nil {
		return reqres.PullRequestReassignResponse{}, err
	}
	if !assigned {
		return reqres.PullRequestReassignResponse{}, dbErrors.ErrorReviewerNotAssigned
	}

	var teamName string
	err = m.Conn.QueryRowContext(ctx, `
		SELECT team_name FROM users WHERE user_id = $1
	`, authorID).Scan(&teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reqres.PullRequestReassignResponse{}, dbErrors.ErrorUserNotFound
		}
		return reqres.PullRequestReassignResponse{}, err
	}

	rows, err := m.Conn.QueryContext(ctx, `
		SELECT user_id FROM users
		WHERE team_name = $1 AND is_active = TRUE AND user_id NOT IN ($2, $3)
	`, teamName, req.OldUserID, authorID)
	if err != nil {
		return reqres.PullRequestReassignResponse{}, err
	}
	defer rows.Close() //nolint:errcheck

	var candidates []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return reqres.PullRequestReassignResponse{}, err
		}
		candidates = append(candidates, uid)
	}

	if len(candidates) == 0 {
		return reqres.PullRequestReassignResponse{}, dbErrors.ErrorNoCandidateForReviewer
	}

	rand.Seed(time.Now().UnixNano())
	newReviewer := candidates[rand.Intn(len(candidates))]

	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		return reqres.PullRequestReassignResponse{}, err
	}
	defer tx.Rollback() //nolint:errcheck

	_, err = tx.ExecContext(ctx, `
		DELETE FROM pr_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2
	`, req.PullRequestID, req.OldUserID)
	if err != nil {
		return reqres.PullRequestReassignResponse{}, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`, req.PullRequestID, newReviewer)
	if err != nil {
		return reqres.PullRequestReassignResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		return reqres.PullRequestReassignResponse{}, err
	}

	var resp reqres.PullRequestReassignResponse
	resp.ReplacedBy = newReviewer
	resp.PR.PullRequestID = req.PullRequestID
	resp.PR.Status = status
	resp.PR.AuthorID = authorID

	rows2, err := m.Conn.QueryContext(ctx, `
		SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1
	`, req.PullRequestID)
	if err != nil {
		return reqres.PullRequestReassignResponse{}, err
	}
	defer rows2.Close() //nolint:errcheck

	for rows2.Next() {
		var rid string
		if err := rows2.Scan(&rid); err != nil {
			return reqres.PullRequestReassignResponse{}, err
		}
		resp.PR.AssignedReviewers = append(resp.PR.AssignedReviewers, rid)
	}

	return resp, nil
}
