package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"team-service/internal"
)

func main() {
	// Connect to MongoDB
	if err := internal.ConnectToMongo(); err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	teamRepo := internal.NewTeamRepository(internal.Client.Database("myblog"))

	// Initialize handler
	handler := internal.NewTeamHandler(teamRepo)

	// Create router
	r := mux.NewRouter()

	// Team routes
	r.HandleFunc("/teams", internal.AuthMiddleware(handler.CreateTeam)).Methods("POST")
	r.HandleFunc("/teams/{id}", handler.GetTeam).Methods("GET")
	r.HandleFunc("/teams/user", internal.AuthMiddleware(handler.GetUserTeams)).Methods("GET")
	r.HandleFunc("/teams/{id}", internal.AuthMiddleware(handler.UpdateTeam)).Methods("PUT")
	r.HandleFunc("/teams/{id}", internal.AuthMiddleware(handler.DeleteTeam)).Methods("DELETE")
	r.HandleFunc("/teams/{id}/invite", internal.AuthMiddleware(handler.InviteMember)).Methods("POST")
	r.HandleFunc("/teams/invites/{id}", internal.AuthMiddleware(handler.RespondToInvite)).Methods("POST")
	r.HandleFunc("/teams/{id}/members/{username}", internal.AuthMiddleware(handler.RemoveMember)).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:          300,
	})

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Default port for team service
	}

	srv := &http.Server{
		Handler:      c.Handler(r),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Team service starting on port %s", port)
	log.Fatal(srv.ListenAndServe())
} 