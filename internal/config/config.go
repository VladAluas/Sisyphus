// Package config exposes the internal configuration to the application
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost string
	DBPort string
	DBName string
	DBUser string
	DBPass string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}
	return Config{
		DBHost: os.Getenv("METADATA_HOST"),
		DBPort: os.Getenv("METADATA_PORT"),
		DBName: os.Getenv("METADATA_NAME"),
		DBUser: os.Getenv("METADATA_USER"),
		DBPass: os.Getenv("METADATA_PASSWORD"),
	}, nil
}
