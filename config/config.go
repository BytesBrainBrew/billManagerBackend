// write function to read .env and connect mongoDB
package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

//function load .env variables
func LoadEnvVariables() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func ConnectDB() *mongo.Client {
	LoadEnvVariables()

    // Read environment variables
    mongoURI := os.Getenv("MONGO_URI")

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(mongoURI))
    if err!= nil {
        fmt.Print("not able to connect to MongoDB: ", err)
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")

    DB = client.Database("Go-User")  // save the client for other functions to use it
    return client
}