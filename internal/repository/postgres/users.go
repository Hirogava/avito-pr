package postgres

import (
	"database/sql"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
)

func (manager *Manager) SetUserIsActive(req reqres.UserSetIsActiveRequest) (reqres.UserResponse, error) {
	var user reqres.UserResponse
	err := manager.Conn.QueryRow(`UPDATE users SET is_active = $1 WHERE user_id = $2 RETURNING is_active, username, team_name, user_id`, req.IsActive, req.UserID).Scan(&user.IsActive, &user.Username, &user.TeamName, &user.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return reqres.UserResponse{}, dbErrors.ErrorUserNotFound
		}
		return reqres.UserResponse{}, err
	}

	return user, nil
}

func (manager *Manager) GetUsersReview(req reqres.UsersGetReviewQuery) (reqres.PullRequestListResponse, error) {
	var reviewList reqres.PullRequestListResponse
	reviewList.UserID = req.UserID

	rows, err := manager.Conn.Query(`
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			pr.author_id,
			pr.status
		FROM pr_reviewers r
		JOIN pull_requests pr ON r.pull_request_id = pr.pull_request_id
		WHERE r.reviewer_id = $1
	`, req.UserID)
	if err != nil {
		return reviewList, err
	}
	defer rows.Close()

	for rows.Next() {
		var pr reqres.PullRequestShortResponse
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return reviewList, err
		}
		reviewList.PullRequests = append(reviewList.PullRequests, pr)
	}

	if err := rows.Err(); err != nil {
		return reviewList, err
	}

	return reviewList, nil
}

func (manager *Manager) GetUsers() ([]reqres.UserResponse, error) {
	rows, err := manager.Conn.Query(`
		SELECT username, team_name, user_id, is_active FROM users
		`)
	if err != nil {
		return nil, err
	}

	var users []reqres.UserResponse
	for rows.Next() {
		var user reqres.UserResponse

		if err := rows.Scan(&user.Username, &user.TeamName, &user.UserID, &user.IsActive); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
