package models

import (
	"context"
	"log"

	"example.com/user-service/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents the user document in MongoDB
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Username string             `bson:"username"`
    Email    string             `bson:"email"`
    Password string             `bson:"password"`
}

func CreateIndexModel(){
    collection := config.DB.Collection("users")

    indexModel := mongo.IndexModel{
        Keys:    bson.M{"email": 1},
        Options: options.Index().SetUnique(true),
    }

    _, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("User index created")
}
