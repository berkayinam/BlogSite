package internal

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func getMongoURI() string {
	// Production/Kubernetes ortamı için
	if os.Getenv("ENVIRONMENT") == "production" {
		return "mongodb://mongo-service:27017"
	}
	
	// Custom MongoDB URI varsa onu kullan
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return uri
	}
	
	// Development ortamı için default
	return "mongodb://localhost:27017"
}

func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := getMongoURI()
	log.Printf("Connecting to MongoDB at: %s\n", mongoURI)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to MongoDB")
	Client = client
}
