// Package prs provides handlers for pull requests
package prs

import (
	"net/http"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/handlers/middleware"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

// InitPRSHandlers - инициализация обработчиков для pull requests
func InitPRSHandlers(r *gin.Engine, manager *postgres.Manager) {
	secureUsers := r.Group("/pullRequest")
	secureUsers.Use(middleware.AuthMiddleware())
	{
		secureUsers.POST("/create", func(c *gin.Context) {
			CreatePR(c, manager)
		})
		secureUsers.POST("/merge", func(c *gin.Context) {
			MergePR(c, manager)
		})
		secureUsers.POST("/reassign", func(c *gin.Context) {
			ReassignAuthor(c, manager)
		})
	}
}

// CreatePR - создание pull request
func CreatePR(c *gin.Context, manager *postgres.Manager) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorTeamNotFound.Error()

		c.JSON(http.StatusForbidden, errResp)
		return
	}

	var req reqres.PullRequestCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := manager.CreatePullRequest(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, gin.H{"pull_request": pr})
	case dbErrors.ErrorUserNotFound:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorUserNotFound.Error()
		c.JSON(http.StatusNotFound, errResp)
	case dbErrors.ErrorPRAlreadyExists:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodePRExists
		errResp.Error.Message = dbErrors.ErrorPRAlreadyExists.Error()
		c.JSON(http.StatusBadRequest, errResp)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// MergePR - слияние pull request
func MergePR(c *gin.Context, manager *postgres.Manager) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorTeamNotFound.Error()

		c.JSON(http.StatusForbidden, errResp)
		return
	}

	var req reqres.PullRequestMergeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := manager.MergePullRequest(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"pull_request": pr})
	case dbErrors.ErrorPRSNotFound:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorPRSNotFound.Error()
		c.JSON(http.StatusNotFound, errResp)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// ReassignAuthor - смена автора pull request
func ReassignAuthor(c *gin.Context, manager *postgres.Manager) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorTeamNotFound.Error()

		c.JSON(http.StatusForbidden, errResp)
		return
	}

	var req reqres.PullRequestReassignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := manager.ReassignPRAuthor(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"pull_request": pr})
	case dbErrors.ErrorUserNotFound:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorUserNotFound.Error()
		c.JSON(http.StatusNotFound, errResp)
	case dbErrors.ErrorNoCandidateForReviewer:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeNoCandidate
		errResp.Error.Message = dbErrors.ErrorNoCandidateForReviewer.Error()
		c.JSON(http.StatusBadRequest, errResp)
	case dbErrors.ErrorPRMerged:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodePRMerged
		errResp.Error.Message = dbErrors.ErrorPRMerged.Error()
		c.JSON(http.StatusBadRequest, errResp)
	case dbErrors.ErrorReviewerNotAssigned:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeNotAssigned
		errResp.Error.Message = dbErrors.ErrorReviewerNotAssigned.Error()
		c.JSON(http.StatusBadRequest, errResp)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
