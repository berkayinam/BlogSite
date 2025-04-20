package internal

import (
	"context"
	"fmt"
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

func ConnectToMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := getMongoURI()
	fmt.Printf("Connecting to MongoDB at: %s\n", mongoURI)

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("mongo connection error: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongo ping error: %v", err)
	}

	fmt.Println("MongoDB connected successfully")
	Client = client
	return nil
}