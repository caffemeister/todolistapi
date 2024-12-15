package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() error {
	var err error

	// Create the db
	db, err = sql.Open("sqlite", "tasks.db")
	if err != nil {
		return fmt.Errorf("failed to connect to db: %v", err)
	}

	// Ping to see if working
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	);`

	// Create the table
	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create tasks table: %v", err)
	}

	return nil
}
