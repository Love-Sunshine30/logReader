package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	dsn := "postgres://postgres:password123@localhost:5432/log_reader?sslmode=disable"

	// open the database
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// ping the db
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping the db: %v", err)
	}
	fmt.Println("Database connected successfully!")
	return nil
}
