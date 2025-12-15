package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	HTTPAddr string
	DB       DBConfig
}

// DBConfig defines database connectivity and pooling settings.
type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// New reads environment variables and applies sane defaults for local development.
func New() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8000"),
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "db"),
			Port:            mustInt(getEnv("DB_PORT", "5432"), 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Name:            getEnv("DB_NAME", "cyber_tournament"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    mustInt(getEnv("DB_MAX_OPEN_CONNS", "25"), 25),
			MaxIdleConns:    mustInt(getEnv("DB_MAX_IDLE_CONNS", "25"), 25),
			ConnMaxLifetime: mustDuration(getEnv("DB_CONN_MAX_LIFETIME", "30m"), 30*time.Minute),
		},
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustInt(raw string, fallback int) int {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}

func mustDuration(raw string, fallback time.Duration) time.Duration {
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	return fallback
}
