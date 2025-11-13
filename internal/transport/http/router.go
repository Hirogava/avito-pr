package http

import (
	"time"

	"github.com/Hirogava/avito-pr/internal/config/logger"
	"github.com/Hirogava/avito-pr/internal/handlers/auth"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	auth.InitAuthHandlers(r, manager)

	logger.Logger.Info("HTTP router created successfully")
	return r
}
