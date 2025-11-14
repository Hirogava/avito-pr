// Package environment provides functionality to load environment variables from a .env file.
package environment

import (
	"github.com/joho/godotenv"
)

// LoadEnvFile - Загружает переменные окружения из файла .env.
func LoadEnvFile(filename string) error {
	return godotenv.Load(filename)
}
