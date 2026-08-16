package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreateUserParams
}

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (c Client) GetUsers() ([]User, error) {
	query := `SELECT id, email, is_verified FROM users;`

	rows, err := c.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		var id string
		if err := rows.Scan(&id, &user.Email, &user.IsVerified); err != nil {
			return nil, err
		}
		user.ID, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (c Client) CreateUser(params CreateUserParams) (*User, error) {
	id := uuid.New()

	query := `
		INSERT INTO users
		    (id, created_at, updated_at, email, password)
		VALUES
		    (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?)
	`
	_, err := c.db.Exec(query, id.String(), params.Email, params.Password)
	if err != nil {
		return nil, err
	}

	return c.GetUser(id)
}

func (c Client) GetUser(id uuid.UUID) (*User, error) {
	query := `
		SELECT id, created_at, updated_at, email, password
		FROM users
		WHERE id = ?
	`
	var user User
	var idStr string
	err := c.db.QueryRow(query, id.String()).Scan(&idStr, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c Client) DeleteUser(id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`
	_, err := c.db.Exec(query, id.String())
	return err
}
