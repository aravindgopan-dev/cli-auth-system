package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL     string
	SessionDuration time.Duration
	LockoutDuration time.Duration
}

func Load() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
	}

	sessDur := 5 * time.Minute
	if val, err := time.ParseDuration(os.Getenv("SESSION_DURATION")); err == nil {
		sessDur = val
	}

	lockDur := 1 * time.Minute
	if val, err := time.ParseDuration(os.Getenv("LOCKOUT_DURATION")); err == nil {
		lockDur = val
	}

	return &Config{
		DatabaseURL:     dbURL,
		SessionDuration: sessDur,
		LockoutDuration: lockDur,
	}
}