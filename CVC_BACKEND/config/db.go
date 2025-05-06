package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// InitDB initializes the MongoDB connection
func InitDB() {

	// Load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatal("❌ Error loading .env file:", err)
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("❌ MONGODB_URI not set in environment")
	}

	clientOptions := options.Client().ApplyURI(uri)

	// Set a timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ Error connecting to MongoDB:", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	// Set global DB variable
	db = client.Database("container_security")
	log.Println("✅ Connected to MongoDB!")
}

// GetDB returns the database instance
func GetDB() *mongo.Database {
	if db == nil {
		log.Fatal("❌ Database is not initialized! Did you call InitDB()?")
	}
	return db
}
