package main

import (
	"log"
	"net/http"
	"auth-service/internal"
)

func main() {
	err := internal.ConnectToMongo()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	
	http.HandleFunc("/health", internal.HealthHandler)
	http.HandleFunc("/login", internal.LoginHandler)
	http.HandleFunc("/register", internal.RegisterHandler)

	log.Println("Auth service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}