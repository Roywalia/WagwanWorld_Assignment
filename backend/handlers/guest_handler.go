// backend/handlers/guest_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Guest struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type GuestHandler struct{ db *sql.DB }

// NewGuestHandler sets up the handler with DB connection
func NewGuestHandler(db *sql.DB) *GuestHandler {
	return &GuestHandler{db}
}

// GetGuests fetches all guests with optional status and search filters
func (h *GuestHandler) GetGuests(w http.ResponseWriter, r *http.Request) {
	// Grab query params
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	query := `
		SELECT id, name, email, phone, status, created_at 
		FROM guests 
		WHERE 1=1
	`
	args := []any{}
	argIndex := 1

	// Add status filter if provided
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Add search filter (name OR email)
	if search != "" {
		query += fmt.Sprintf(" AND (LOWER(name) LIKE LOWER($%d) OR LOWER(email) LIKE LOWER($%d))", argIndex, argIndex+1)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	// Run the query
	rows, err := h.db.Query(query, args...)
	if err != nil {
		sendJSONError(w, 500, "Failed to load guests")
		return
	}
	defer rows.Close()

	var guests []Guest
	for rows.Next() {
		var g Guest
		// Scan into struct fields
		if err := rows.Scan(&g.ID, &g.Name, &g.Email, &g.Phone, &g.Status, &g.CreatedAt); err != nil {
			sendJSONError(w, 500, "Failed to scan guest")
			return
		}
		guests = append(guests, g)
	}

	sendJSON(w, 200, guests)
}

// CreateGuest adds a new guest via the admin panel
func (h *GuestHandler) CreateGuest(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Status string `json:"status"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendJSONError(w, 400, "Invalid JSON")
		return
	}

	// Basic validation
	if input.Name == "" || input.Email == "" {
		sendJSONError(w, 400, "Name and email required")
		return
	}
	if !isValidEmail(input.Email) {
		sendJSONError(w, 400, "Invalid email format")
		return
	}

	// Default to pending, allow attending/declined
	status := "pending"
	if input.Status != "" {
		switch input.Status {
		case "attending", "declined":
			status = input.Status
		}
	}

	// Insert and get the new ID
	var id int
	err := h.db.QueryRow(`
		INSERT INTO guests (name, email, phone, status) 
		VALUES ($1, $2, $3, $4) RETURNING id
	`, input.Name, input.Email, input.Phone, status).Scan(&id)
	if err != nil {
		// Log this in prod
		sendJSONError(w, 500, "Failed to create guest")
		return
	}

	sendJSON(w, 201, Guest{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		Status:    status,
		CreatedAt: time.Now().Format(time.RFC3339), // Fresh timestamp
	})
}

// DeleteGuest removes a guest by ID
func (h *GuestHandler) DeleteGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Simple delete
	_, err := h.db.Exec("DELETE FROM guests WHERE id = $1", id)
	if err != nil {
		sendJSONError(w, 500, "Failed to delete guest")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// === Helper Functions ===

// sendJSON writes JSON with proper headers
func sendJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// sendJSONError quick error response
func sendJSONError(w http.ResponseWriter, code int, message string) {
	sendJSON(w, code, map[string]string{"error": message})
}

// isValidEmail does a light check — good enough for this app
func isValidEmail(email string) bool {
	if len(email) < 6 {
		return false
	}
	at, dot := -1, -1
	for i, ch := range email {
		if ch == '@' {
			if at != -1 {
				return false
			}
			at = i
		}
		if ch == '.' && at != -1 && i > at+1 {
			dot = i
		}
	}
	return at > 0 && dot > at+1 && dot < len(email)-1
}

// convert empty string → NULL for DB
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
