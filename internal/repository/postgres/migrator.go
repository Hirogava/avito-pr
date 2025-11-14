// Package postgres implements the repository interface for PostgreSQL.
package postgres

import (
	"fmt"

	"github.com/Hirogava/avito-pr/internal/config/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	// import file source for migrations.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrate - миграция БД
func (manager *Manager) Migrate() {
	logger.Logger.Debug("Starting database migrations")

	driver, err := postgres.WithInstance(manager.Conn, &postgres.Config{})
	if err != nil {
		logger.Logger.Fatal("Failed to create migration driver", "error", err.Error())
		panic(fmt.Sprintf("Не удалось создать драйвер миграции: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/internal/repository/postgres/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		logger.Logger.Fatal("Failed to create migrator", "error", err.Error())
		panic(fmt.Sprintf("Не удалось создать мигратора: %v", err))
	}

	logger.Logger.Debug("Running database migrations")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Logger.Fatal("Failed to apply migrations", "error", err.Error())
		panic(fmt.Sprintf("Не удалось применить миграции: %v", err))
	}

	if err == migrate.ErrNoChange {
		logger.Logger.Info("No new migrations to apply")
	} else {
		logger.Logger.Info("Database migrations completed successfully")
	}
}
