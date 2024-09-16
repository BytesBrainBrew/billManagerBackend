package models

import (
	"context"
	"log"

	"example.com/user-service/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bill struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	BillTitle string             `bson:"bill_title"`
	Amount    float64            `bson:"amount"`
}


func CreateBillIndexModel(){
	collection := config.DB.Collection("bills")

    indexModel := mongo.IndexModel{
        Keys:    bson.M{"email": 1},
    }

    _, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
    if err!= nil {
        log.Fatal(err)
    }
    log.Println("Bill index created")
}
