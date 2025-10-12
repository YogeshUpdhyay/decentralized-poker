package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// DatabaseConfig holds configuration for initializing the database
type DatabaseConfig struct {
	Path string // Path to the SQLite database file
}

// Database struct holds the DB connection
type Database struct {
	conn *sql.DB
}

var db *Database

// InitDatabase initializes and connects to the SQLite database
func InitDatabase(cfg DatabaseConfig) error {
	conn, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	// Optionally, ping to check connection
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db = &Database{conn: conn}

	return nil
}

// Get returns the database connection
func Get() *Database {
	return db
}

// BootstrapTables creates required tables if they do not exist
func (d *Database) BootstrapTables() error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS user_metadata (
        username VARCHAR(100) PRIMARY KEY,
        last_login_ts INTEGER,
        create_ts INTEGER,
        update_ts INTEGER
    );`
	_, err := d.conn.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create user_metadata table: %w", err)
	}
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
