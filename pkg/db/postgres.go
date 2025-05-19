package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"wallet-go/pkg/config"
)

func InitPostgres() *sql.DB {
	config.LoadEnv()
	dbCfg := config.GetDBConfig()

	// Connect to default "postgres" database to create wallet DB if missing
	defaultDSN := dbCfg.GetDefaultDSN()
	defaultDB, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		log.Fatalf("failed to connect to default DB: %v", err)
	}
	defer defaultDB.Close()

	// Create wallet DB
	if err := EnsureDatabaseExists(defaultDB, dbCfg.Name); err != nil {
		log.Fatalf("DB creation error: %v", err)
	}

	mainDSN := dbCfg.GetMainDSN()
	schemaPath := "pkg/db/schema/schema.sql"

	mainDBConn, err := sql.Open("postgres", mainDSN)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	if err := mainDBConn.Ping(); err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	if err := InitSchema(mainDBConn, schemaPath); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	return mainDBConn
}
