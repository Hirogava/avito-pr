package team

import (
	"net/http"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/handlers/middleware"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

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
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamAlreadyExists
		error.Error.Message = dbErrors.ErrorTeamAlreadyExists.Error()
		c.JSON(http.StatusBadRequest, error)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

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
		var error reqres.ErrorResponse
		error.Error.Code = dbErrors.CodeTeamNotFound
		error.Error.Message = dbErrors.ErrorTeamNotFound.Error()
		c.JSON(http.StatusNotFound, error)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
