package main

import (
	"database/sql"
	"event-guest-manager/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "123456")
	dbName := getEnv("DB_NAME", "eventguests")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Wait for database to be ready
	for i := 0; i < 30; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Println("Waiting for database...")
		time.Sleep(time.Second)
	}

	if err != nil {
		log.Fatal("Database not available:", err)
	}

	log.Println("Successfully connected to database")

	// Initialize handlers
	guestHandler := handlers.NewGuestHandler(db)
	eventHandler := handlers.NewEventHandler(db)

	// Setup router
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	// Events
	api.HandleFunc("/events", eventHandler.GetEvents).Methods("GET")
	api.HandleFunc("/events/{id}/rsvp", eventHandler.CreateRSVP).Methods("POST")
	
	// API routes
	
	api.HandleFunc("/guests", guestHandler.GetGuests).Methods("GET")
	api.HandleFunc("/guests", guestHandler.CreateGuest).Methods("POST")
	api.HandleFunc("/guests/{id}", guestHandler.GetGuests).Methods("GET")
	api.HandleFunc("/guests/{id}", guestHandler.DeleteGuest).Methods("DELETE")

	api.HandleFunc("/events", eventHandler.CreateEvent).Methods("POST")

	// CORS setup
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:4173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
