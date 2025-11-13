package prs

import (
	"net/http"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/handlers/middleware"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

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

func CreatePR(c *gin.Context, manager *postgres.Manager) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorTeamNotFound.Error()

		c.JSON(http.StatusForbidden, error)
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
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorUserNotFound.Error()
		c.JSON(http.StatusNotFound, error)
	case dbErrors.ErrorPRAlreadyExists:
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodePRExists
		error.Error.Message = dbErrors.ErrorPRAlreadyExists.Error()
		c.JSON(http.StatusBadRequest, error)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func MergePR(c *gin.Context, manager *postgres.Manager) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorTeamNotFound.Error()

		c.JSON(http.StatusForbidden, error)
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
		c.JSON(http.StatusCreated, gin.H{"pull_request": pr})
	case dbErrors.ErrorPRSNotFound:
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorPRSNotFound.Error()
		c.JSON(http.StatusNotFound, error)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func ReassignAuthor(c *gin.Context, manager *postgres.Manager) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorTeamNotFound.Error()

		c.JSON(http.StatusForbidden, error)
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
		c.JSON(http.StatusCreated, gin.H{"pull_request": pr})
	case dbErrors.ErrorUserNotFound:
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorUserNotFound.Error()
		c.JSON(http.StatusNotFound, error)
	case dbErrors.ErrorNoCandidateForReviewer:
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeNoCandidate
		error.Error.Message = dbErrors.ErrorNoCandidateForReviewer.Error()
		c.JSON(http.StatusBadRequest, error)
	case dbErrors.ErrorPRMerged:
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodePRMerged
		error.Error.Message = dbErrors.ErrorPRMerged.Error()
		c.JSON(http.StatusBadRequest, error)
	case dbErrors.ErrorReviewerNotAssigned:
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeNotAssigned
		error.Error.Message = dbErrors.ErrorReviewerNotAssigned.Error()
		c.JSON(http.StatusBadRequest, error)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
