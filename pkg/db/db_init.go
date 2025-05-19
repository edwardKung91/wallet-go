package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// EnsureDatabaseExists connects to the default DB and creates your wallet DB if needed.
func EnsureDatabaseExists(defaultDB *sql.DB, dbName string) error {
	query := fmt.Sprintf(`SELECT 1 FROM pg_database WHERE datname = '%s'`, dbName)
	var exists int
	err := defaultDB.QueryRow(query).Scan(&exists)
	if err == sql.ErrNoRows {
		_, err := defaultDB.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			return fmt.Errorf("failed to create database %s: %w", dbName, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking for database existence: %w", err)
	}
	return nil
}

// InitSchema runs the schema SQL from the given file on the connected DB.
func InitSchema(db *sql.DB, schemaPath string) error {
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	queries := strings.TrimSpace(string(content))
	_, err = db.Exec(queries)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}
