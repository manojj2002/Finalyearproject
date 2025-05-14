package main

import (
	"CVC_ragh/config"
	"CVC_ragh/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},        // Allow frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	// Register user routes
	config.InitDB()
	routes.RegisterUserRoutes(r)
	routes.RegisterContainerRoutes(r)
	routes.RegisterScanningRoutes(r)
	routes.RegisterDynamicScanningRoutes(r)
	routes.RegisterFalcoWebHook(r)
	routes.RegisterGitHubAuthRoutes(r)
	//routes.RegisterDashboardRoutes(r)
	r.Run(":4000") // Start server on port 5000
}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func migrateData() {
// 	// Connect to local MongoDB
// 	localClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		log.Fatal("Error connecting to local MongoDB:", err)
// 	}

// 	// Connect to MongoDB Atlas
// 	cloudClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://jrraghav:bmsce2025@cluster0.jgir2hv.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
// 	if err != nil {
// 		log.Fatal("Error connecting to MongoDB Atlas:", err)
// 	}

// 	// Fetch data from local MongoDB
// 	localCollection := localClient.Database("container_security").Collection("trivy_scans")
// 	cursor, err := localCollection.Find(context.TODO(), bson.M{})
// 	if err != nil {
// 		log.Fatal("Error fetching data from local MongoDB:", err)
// 	}

// 	// Insert data into MongoDB Atlas
// 	cloudCollection := cloudClient.Database("container_security").Collection("trivy_scans")
// 	for cursor.Next(context.TODO()) {
// 		var document bson.M
// 		if err := cursor.Decode(&document); err != nil {
// 			log.Fatal("Error decoding document:", err)
// 		}

// 		// Insert into Atlas
// 		_, err = cloudCollection.InsertOne(context.TODO(), document)
// 		if err != nil {
// 			log.Fatal("Error inserting document into Atlas:", err)
// 		}
// 	}

// 	fmt.Println("Migration successful!")
// }

// func main() {
// 	migrateData()
// }
