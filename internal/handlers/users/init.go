package users

import (
	"net/http"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/handlers/middleware"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

func InitUsersHandlers(r *gin.Engine, manager *postgres.Manager) {
	secureUsers := r.Group("/team")
	secureUsers.Use(middleware.AuthMiddleware())
	{
		secureUsers.POST("/setIsActive", func(c *gin.Context) {
			SetIsActive(c, manager)
		})
		secureUsers.GET("/getReview", func(c *gin.Context) {
			GetReview(c, manager)
		})
	}
}

func SetIsActive(c *gin.Context, manager *postgres.Manager) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorTeamNotFound.Error()

		c.JSON(http.StatusForbidden, error)
		return
	}

	var req reqres.UserSetIsActiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := manager.SetUserIsActive(req)
	switch err {
		case nil:
			c.JSON(http.StatusOK, gin.H{"user": user})
		case dbErrors.ErrorUserNotFound:
			var error reqres.ErrorResponse
			error.Error.Code = dbErrors.CodeTeamNotFound
			error.Error.Message = dbErrors.ErrorTeamNotFound.Error()
			c.JSON(http.StatusBadRequest, error)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetReview(c *gin.Context, manager *postgres.Manager) {
	var req reqres.UsersGetReviewQuery

	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prs, err := manager.GetUsersReview(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, prs)
	case dbErrors.ErrorPRSNotFound:
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorTeamNotFound.Error()
		c.JSON(http.StatusBadRequest, error)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
