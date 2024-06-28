// Package db : Database package
package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

// NewMySQLStorage creates a new MySQL storage instance
func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(151)
	db.SetMaxIdleConns(151)

	return db, nil
}
