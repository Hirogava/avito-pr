// Package reqres models for responses and requests
package reqres

// TeamGetQuery - Query параметры для /team/get.
type TeamGetQuery struct {
	TeamName string `form:"team_name" binding:"required"`
}

// UsersGetReviewQuery - Query параметры для /users/getReview.
type UsersGetReviewQuery struct {
	UserID string `form:"user_id" binding:"required"`
}
