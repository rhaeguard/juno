package db

import (
	"database/sql"
	"fmt"
	"github.com/mensurowary/juno/config"
	log "github.com/sirupsen/logrus"
)

func Initialize() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.DatabaseConfig.Username,
		config.DatabaseConfig.Password,
		config.DatabaseConfig.Endpoint,
		config.DatabaseConfig.Database,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
