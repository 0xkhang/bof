package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Birthday struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreateBirthdayParams
}

type CreateBirthdayParams struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	DOB    time.Time `json:"dob"`
}

func (c Client) GetBirthdayByName(userID uuid.UUID, name string) (Birthday, error) {
	query := `
        SELECT id, name, dob FROM birthdays
        WHERE name = ? and user_id = ?
    	`

	var b Birthday
	err := c.db.QueryRow(query, name, userID).Scan(&b.ID, &b.Name, &b.DOB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Birthday{}, err
		}
	}

	return b, nil
}

func (c Client) CreateBirthday(params CreateBirthdayParams) (Birthday, error) {
	id := uuid.New()

	query := `
        INSERT INTO birthdays 
            (id, user_id, name, dob, created_at, updated_at)
        VALUES
            (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `

	_, err := c.db.Exec(query, id.String(), params.UserID, params.Name, params.DOB)
	if err != nil {
		return Birthday{}, err
	}

	return c.GetBirthdayByName(params.UserID, params.Name)
}
