package database

import (
	"context"
	"database/sql"
	"log"
)

type Client struct {
	db *sql.DB
}

func NewClient(pathToDB string) (Client, error) {
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		return Client{}, err
	}

	if db.Ping() == nil {
		log.Printf("DB works!")
	}

	_, err = db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;")
	if err != nil {
		return Client{}, err
	}

	c := Client{db}
	err = c.autoMigrate()
	if err != nil {
		return Client{}, err
	}
	return c, nil
}

func (c *Client) autoMigrate() error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		is_verified BOOLEAN DEFAULT 0,
		token_verification_hash TEXT NULL,
		token_expires_at TIMESTAMP NULL,
		hashed_password TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);
	`

	_, err := c.db.ExecContext(context.Background(), userTable)
	if err != nil {
		return err
	}

	birthdayTable := `
	CREATE TABLE IF NOT EXISTS birthdays (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		friends_name TEXT NOT NULL,
		dob TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE(user_id, friends_name)
	);
	`

	_, err = c.db.ExecContext(context.Background(), birthdayTable)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Reset() error {
	usersTable := `
		DELETE FROM users;
	`

	_, err := c.db.ExecContext(context.Background(), usersTable)
	if err != nil {
		return err
	}
	return nil
}
