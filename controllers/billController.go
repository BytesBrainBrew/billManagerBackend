package controllers

import (
	"strconv"

	"example.com/user-service/config"
	"example.com/user-service/models"
	"example.com/user-service/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getBillCollection() *mongo.Collection {
	return config.DB.Collection("bills")
}

func AddBill(c *fiber.Ctx) error {
	var bill models.Bill
	billCollection := getBillCollection()
	userCollection := getUserCollection()

	if err := c.BodyParser(&bill); err!= nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Error parsing request body",
        })
    }

	tokenStr := c.Get("Authorization")
	email, isExist := utils.IsTokenExpired(tokenStr)
	if!isExist {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Token not found or expired",
        })
    }
	bill.Email = email
	var existUser models.User
	filter := bson.M{"email": email}
	err := userCollection.FindOne(c.Context(), filter).Decode(&existUser)
	if err!= nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error finding user in database",
        })
    }
	if existUser.Email != email {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "User not found",
        })
	}

	result, err := billCollection.InsertOne(c.Context(), bill)
	if err!= nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error inserting bill into database",
        })
    }
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Bill added successfully",
        "billId":  result.InsertedID,
    })
}

//add pagination logic here
func GetBills(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	email, isExist := utils.IsTokenExpired(tokenStr)
	if!isExist {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Token not found or expired",
        })
    }
	
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err!= nil || page < 1 {
        page = 1
    }

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err!= nil || limit < 1 {
        limit = 10
    }

	skip := (page - 1) * limit

	var bills []models.Bill
    billCollection := getBillCollection()
    filter := bson.M{"email": email}

	options := options.Find()
    options.SetSkip(int64(skip))
    options.SetLimit(int64(limit))

    cursor, err := billCollection.Find(c.Context(), filter, options)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error fetching bills from database",
        })
    }
    defer cursor.Close(c.Context())

    for cursor.Next(c.Context()) {
        var bill models.Bill
        err := cursor.Decode(&bill)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Error decoding bill",
            })
        }
        bills = append(bills, bill)
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Bills fetched successfully",
        "data":    bills,
        "page":    page,
        "limit":   limit,
        "total":   len(bills),
    })

}