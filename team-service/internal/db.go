package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func ConnectPostgres(host, port, user, password, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("error creating tables: %v", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	// Teams table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Team members table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS team_members (
			team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
			username VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (team_id, username)
		)
	`)
	if err != nil {
		return err
	}

	// Team invites table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS team_invites (
			id SERIAL PRIMARY KEY,
			team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
			username VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	return nil
} 