package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// LoadEnv loads the .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, relying on environment variables")
	}
}

// GetDBConfig returns the database configuration
func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

// GetMainDSN builds the PostgreSQL DSN string for the main DB
func (cfg DBConfig) GetMainDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)
}

// GetDefaultDSN builds the PostgreSQL DSN string for the default user
func (cfg DBConfig) GetDefaultDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port,
	)
}
