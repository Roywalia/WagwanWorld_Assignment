package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Event represents a single event with nullable fields for flexible JSON output
type Event struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description,omitempty"`
	EventDate   sql.NullTime   `json:"event_date"`
	Location    sql.NullString `json:"location,omitempty"`
	CreatedAt   sql.NullTime   `json:"created_at,omitempty"`
	Display     string         `json:"display"`
	RSVPs       int            `json:"rsvps"`
}

type EventHandler struct{ db *sql.DB }

// NewEventHandler creates a new handler with DB connection
func NewEventHandler(db *sql.DB) *EventHandler {
	return &EventHandler{db}
}

// Custom JSON marshaling to convert sql.NullTime → ISO string (or omit if null)
func (e Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	return json.Marshal(&struct {
		*Alias
		EventDate string `json:"event_date,omitempty"`
		CreatedAt string `json:"created_at,omitempty"`
	}{
		Alias:     (*Alias)(&e),
		EventDate: nullTimeToString(e.EventDate),
		CreatedAt: nullTimeToString(e.CreatedAt),
	})
}

// GetEvents returns all events with a nice display string and RSVP count
func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Fetch raw event data
	rows, err := h.db.Query("SELECT id, title, description, event_date, location, created_at FROM events")
	if err != nil {
		sendJSONError(w, 500, "Failed to load events")
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		// Scan into nullable types
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.EventDate,
			&e.Location,
			&e.CreatedAt,
		); err != nil {
			log.Printf("Scan error: %v", err)
			sendJSONError(w, 500, "Failed to scan event")
			return
		}

		// Convert nullable fields to strings
		desc := nullStringToString(e.Description)
		loc := nullStringToString(e.Location)
		date := ""
		if e.EventDate.Valid {
			date = e.EventDate.Time.Format("02 Jan 06")
		}

		// Truncate long text to avoid dropdown overflow
		title := truncate(e.Title, 30)
		descShort := truncate(desc, 50)
		locShort := truncate(loc, 30)

		parts := []string{title}
		if descShort != "" {
			parts = append(parts, descShort)
		}
		if locShort != "" {
			parts = append(parts, locShort)
		}
		if date != "" {
			parts = append(parts, date)
		}
		e.Display = strings.Join(parts, " – ")

		// Count total RSVPs for this event
		var rsvpCount int
		err = h.db.QueryRow("SELECT COUNT(*) FROM guests WHERE event_id = $1", e.ID).Scan(&rsvpCount)
		if err != nil {
			log.Printf("RSVP count error for event %d: %v", e.ID, err)
			rsvpCount = 0
		}
		e.RSVPs = rsvpCount

		events = append(events, e)
	}

	// Send final JSON
	sendJSON(w, 200, events)
}

// CreateRSVP handles RSVP submission with duplicate check
func (h *EventHandler) CreateRSVP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSONError(w, 400, "Invalid event ID")
		return
	}

	// Parse request body
	var input struct {
		Name                string `json:"name"`
		Email               string `json:"email"`
		Phone               string `json:"phone"`
		RSVPStatus          string `json:"rsvp_status"`
		Notes               string `json:"notes"`
		PlusOnes            int    `json:"plus_ones"`
		DietaryRestrictions string `json:"dietary_restrictions"`
	}
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

	// Map frontend status → DB status
	status := "pending"
	switch input.RSVPStatus {
	case "attending":
		status = "attending"
	case "maybe":
		status = "pending"
	case "declined":
		status = "declined"
	}

	// Check for duplicate RSVP (same email + event)
	var exists bool
	err = h.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM guests WHERE event_id = $1 AND email = $2)
	`, eventID, input.Email).Scan(&exists)

	// Fallback: if table changed, just check email
	if err != nil {
		if strings.Contains(err.Error(), `column "event_id"`) {
			err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM guests WHERE email = $1)", input.Email).Scan(&exists)
		}
		if err != nil {
			sendJSONError(w, 500, "Database error")
			return
		}
	}
	if exists {
		sendJSONError(w, 409, "You already RSVP'd to this event")
		return
	}

	// Insert new guest
	_, err = h.db.Exec(`
		INSERT INTO guests (event_id, name, email, phone, status, notes, plus_ones, dietary_restrictions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, eventID, input.Name, input.Email, input.Phone, status, input.Notes, input.PlusOnes, input.DietaryRestrictions)
	if err != nil {
		log.Printf("Insert failed: %v", err)
		sendJSONError(w, 500, "Failed to save RSVP")
		return
	}

	// Success!
	sendJSON(w, 201, map[string]string{"message": "RSVP saved!"})
}

// CreateEvent adds a new event from the admin panel
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
		EventDate   string `json:"event_date"`
		Location    string `json:"location,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendJSONError(w, 400, "Invalid JSON")
		return
	}

	log.Println(r.Body)

	if input.Title == "" {
		sendJSONError(w, 400, "Title is required")
		return
	}
	if input.EventDate == "" {
		sendJSONError(w, 400, "Event date is required")
		return
	}

	parsedDate, err := time.Parse(time.RFC3339, input.EventDate)
	if err != nil {
		sendJSONError(w, 400, "Invalid date format – use ISO (e.g., 2025-08-20T18:00:00Z)")
		return
	}

	var id int
	err = h.db.QueryRow(`
		INSERT INTO events (title, description, event_date, location)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, input.Title, nullString(input.Description), parsedDate, nullString(input.Location)).Scan(&id)
	if err != nil {
		log.Printf("Failed to create event: %v", err)
		sendJSONError(w, 500, "Failed to create event")
		return
	}

	sendJSON(w, 201, map[string]any{
		"id":          id,
		"title":       input.Title,
		"description": input.Description,
		"event_date":  parsedDate.Format(time.RFC3339),
		"location":    input.Location,
		"message":     "Event created!",
	})
}

// === Helper Functions ===

// truncate cuts long strings with "..." (used in dropdown display)
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// nullStringToString safely converts sql.NullString → string
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// nullTimeToString converts sql.NullTime → ISO string (for JSON)
func nullTimeToString(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
}
