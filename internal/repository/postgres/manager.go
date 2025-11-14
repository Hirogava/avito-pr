// Package postgres implements the repository interface for PostgreSQL.
package postgres

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Hirogava/avito-pr/internal/config/logger"
)

// Manager - DB Manager
type Manager struct {
	Conn *sql.DB
	WG   *sync.WaitGroup
	MU   *sync.RWMutex
}

// NewManager - создание менеджера БД
func NewManager(driverName string, sourceName string) *Manager {
	logger.Logger.Debug("Opening database connection", "driver", driverName)

	db, err := sql.Open(driverName, sourceName)
	if err != nil {
		logger.Logger.Fatal("Failed to open database connection", "error", err.Error())
		panic(fmt.Sprintf("couldn't connect to the database: %v", err))
	}

	logger.Logger.Debug("Pinging database to verify connection")
	if err = db.Ping(); err != nil {
		logger.Logger.Fatal("Database ping failed", "error", err.Error())
		panic(fmt.Sprintf("the database is not responding: %v", err))
	}

	logger.Logger.Info("Database connection established successfully")

	return &Manager{
		Conn: db,
		WG:   &sync.WaitGroup{},
		MU:   &sync.RWMutex{},
	}
}

// Close - закрытие соединения с БД
func (manager *Manager) Close() {
	if manager.Conn != nil {
		logger.Logger.Info("Closing database connection")
		manager.Conn.Close() //nolint:errcheck
		manager.Conn = nil
		logger.Logger.Info("Database connection closed")
	}
}
