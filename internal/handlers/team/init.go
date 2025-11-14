// Package team provides handlers for team
package team

import (
	"net/http"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/handlers/middleware"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

// InitTeamHandlers - инициализация обработчиков для team
func InitTeamHandlers(r *gin.Engine, manager *postgres.Manager) {
	team := r.Group("/team")
	{
		team.POST("/add", func(c *gin.Context) {
			CreateTeam(c, manager)
		})
	}

	secureTeam := r.Group("/team")
	secureTeam.Use(middleware.AuthMiddleware())
	{
		secureTeam.GET("/get", func(c *gin.Context) {
			GetTeam(c, manager)
		})
	}
}

// CreateTeam - создание команды
func CreateTeam(c *gin.Context, manager *postgres.Manager) {
	var req reqres.TeamAddRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := manager.CreateTeam(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, gin.H{"team": team})
	case dbErrors.ErrorTeamAlreadyExists:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamAlreadyExists
		errResp.Error.Message = dbErrors.ErrorTeamAlreadyExists.Error()
		c.JSON(http.StatusBadRequest, errResp)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// GetTeam - получение команды
func GetTeam(c *gin.Context, manager *postgres.Manager) {
	var req reqres.TeamGetQuery

	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := manager.GetTeam(req.TeamName)
	switch err {
	case nil:
		c.JSON(http.StatusOK, team)
	case dbErrors.ErrorTeamNotFound:
		var errResp reqres.ErrorResponse
		errResp.Error.Code = dbErrors.CodeTeamNotFound
		errResp.Error.Message = dbErrors.ErrorTeamNotFound.Error()
		c.JSON(http.StatusNotFound, errResp)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
