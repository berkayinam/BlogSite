package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"team-service/internal"
)

func main() {
	// Get PostgreSQL connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// Create repository and handler
	repo := internal.NewTeamRepository(db)
	handler := internal.NewTeamHandler(repo)

	// Create router
	r := mux.NewRouter()

	// Team endpoints
	r.HandleFunc("/teams", internal.AuthMiddleware(handler.CreateTeam)).Methods("POST")
	r.HandleFunc("/teams/{id}", internal.AuthMiddleware(handler.GetTeam)).Methods("GET")
	r.HandleFunc("/teams/user", internal.AuthMiddleware(handler.GetUserTeams)).Methods("GET")
	r.HandleFunc("/teams/{id}", internal.AuthMiddleware(handler.UpdateTeam)).Methods("PUT")
	r.HandleFunc("/teams/{id}", internal.AuthMiddleware(handler.DeleteTeam)).Methods("DELETE")

	// Team member endpoints
	r.HandleFunc("/teams/invite", internal.AuthMiddleware(handler.InviteMember)).Methods("POST")
	r.HandleFunc("/teams/invite/respond", internal.AuthMiddleware(handler.RespondToInvite)).Methods("POST")
	r.HandleFunc("/teams/members/{teamId}/{username}", internal.AuthMiddleware(handler.RemoveMember)).Methods("DELETE")
	r.HandleFunc("/teams/invites", internal.AuthMiddleware(handler.GetUserInvites)).Methods("GET")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("Team service starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 