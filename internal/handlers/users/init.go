// Package users provides handlers for users
package users

import (
	"net/http"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/handlers/middleware"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

// InitUsersHandlers - инициализация обработчиков для users
func InitUsersHandlers(r *gin.Engine, manager *postgres.Manager) {
	users := r.Group("/users")
	{
		users.GET("", func(c *gin.Context) {
			GerUsers(c, manager)
		})
	}

	secureUsers := r.Group("/users")
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

// GerUsers - получение всех пользователей
func GerUsers(c *gin.Context, manager *postgres.Manager) {
	users, err := manager.GetUsers()
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"users": users})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// SetIsActive - изменение статуса пользователя
func SetIsActive(c *gin.Context, manager *postgres.Manager) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorTeamNotFound.Error()

		c.JSON(http.StatusForbidden, errResp)
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
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorTeamNotFound.Error()
		c.JSON(http.StatusNotFound, errResp)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// GetReview - получение всех pull request
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
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorTeamNotFound.Error()
		c.JSON(http.StatusNotFound, errResp)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
