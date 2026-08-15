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
		password TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);
	`
	_, err := c.db.ExecContext(context.Background(), userTable)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Rollback() error {
	usersTable := `
		DROP TABLE users;
	`

	_, err := c.db.ExecContext(context.Background(), usersTable)
	if err != nil {
		return err
	}
	return nil
}
