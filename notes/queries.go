package db

import (
	"database/sql"
	"errors"
	"time"
)

type Guest struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GetAllGuests retrieves all guests from the database
// BUG #1: This query sorts by created_at instead of name
func GetAllGuests(db *sql.DB, status string) ([]Guest, error) {
	var rows *sql.Rows
	var err error

	if status != "" {
		// BUG #3 (Backend part): This expects 'status' parameter
		rows, err = db.Query(`
			SELECT id, name, email, phone, status, created_at 
			FROM guests 
			WHERE status = $1 
			ORDER BY created_at ASC
		`, status)
	} else {
		rows, err = db.Query(`
			SELECT id, name, email, phone, status, created_at 
			FROM guests 
			ORDER BY created_at ASC
		`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guests []Guest
	for rows.Next() {
		var g Guest
		err := rows.Scan(&g.ID, &g.Name, &g.Email, &g.Phone, &g.Status, &g.CreatedAt)
		if err != nil {
			return nil, err
		}
		guests = append(guests, g)
	}

	return guests, rows.Err()
}

// GetGuestByID retrieves a single guest by ID
func GetGuestByID(db *sql.DB, id int) (*Guest, error) {
	var g Guest
	err := db.QueryRow(`
		SELECT id, name, email, phone, status, created_at 
		FROM guests 
		WHERE id = $1
	`, id).Scan(&g.ID, &g.Name, &g.Email, &g.Phone, &g.Status, &g.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("guest not found")
	}
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// CreateGuest inserts a new guest into the database
func CreateGuest(db *sql.DB, name, email, phone, status string) (*Guest, error) {
	var g Guest
	err := db.QueryRow(`
		INSERT INTO guests (name, email, phone, status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, name, email, phone, status, created_at
	`, name, email, phone, status).Scan(&g.ID, &g.Name, &g.Email, &g.Phone, &g.Status, &g.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &g, nil
}

// DeleteGuest removes a guest from the database
func DeleteGuest(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM guests WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("guest not found")
	}

	return nil
}
