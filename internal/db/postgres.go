package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"wallet-go/internal/config"
)

func InitPostgres() *sql.DB {
	config.LoadEnv()
	dbCfg := config.GetDBConfig()
	dsn := dbCfg.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	return db
}
