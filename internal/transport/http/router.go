package http

import (
	"time"

	"github.com/Hirogava/avito-pr/internal/config/logger"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func CreateRouter(manager *postgres.Manager) *gin.Engine {
	logger.Logger.Debug("Creating HTTP router")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	logger.Logger.Debug("Registering game handlers")
	// game.InitGameHandlers(r, manager) пример инита роутов

	logger.Logger.Info("HTTP router created successfully")
	return r
}
