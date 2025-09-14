package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "secret"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_NAME", "food-app"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c *DBConfig) GetPgConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
