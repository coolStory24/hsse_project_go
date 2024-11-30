package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	Connection *sql.DB
}

func NewDatabase() (*Database, error) {
	url := os.Getenv("DB_URL")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	if url == "" || username == "" || password == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s?sslmode=disable", username, password, url)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{Connection: db}, nil
}
