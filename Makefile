.PHONY: build run clean docker-build docker-up docker-down migrate-up migrate-down

-include .env
export

APP_NAME := avito-pr
CMD_PATH := ./cmd/pr
BUILD_PATH := ./bin
MIGRATIONS_PATH := ./internal/repository/postgres/migrations
GO_VERSION := 1.24.0

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_PATH)/$(APP_NAME) $(CMD_PATH)

run: build
	@echo "Starting $(APP_NAME)..."
	@$(BUILD_PATH)/$(APP_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_PATH)
	@go clean

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

docker-down:
	@echo "Stopping and removing services with Docker Compose..."
	@docker-compose down

migrate-up:
	@echo "Applying database migrations..."
	@docker-compose run --rm app sh -c 'migrate -path /app/internal/repository/postgres/migrations -database "postgres://$${DB_USER}:$${DB_PASSWORD}@db:5432/$${DB_NAME}?sslmode=disable" up'

migrate-down:
	@echo "Reverting last database migration..."
	@docker-compose run --rm app sh -c 'migrate -path /app/internal/repository/postgres/migrations -database "postgres://$${DB_USER}:$${DB_PASSWORD}@db:5432/$${DB_NAME}?sslmode=disable" down 1'