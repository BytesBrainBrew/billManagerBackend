package controllers

import (
	"log"

	"example.com/user-service/config"
	"example.com/user-service/models"
	"example.com/user-service/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return config.DB.Collection("users")
}

func Signup(c *fiber.Ctx) error {
	var user models.User
	userCollection := getUserCollection()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
		})
	}

	var existUser models.User
	filter := bson.M{"email": user.Email}
	err := userCollection.FindOne(c.Context(), filter).Decode(&existUser)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "message": "Email already exists",
        })
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}
	user.Password = hashedPassword

	result, err := userCollection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error inserting user into database",
		})
	}
	log.Println(result.InsertedID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"userId":  result.InsertedID,
	})
}

func Login(c *fiber.Ctx) error {
	var user models.User
    userCollection := getUserCollection()

    if err := c.BodyParser(&user); err!= nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Error parsing request body",
        })
    }

	var updateUser models.User
    filter := bson.M{"email": user.Email}
    err := userCollection.FindOne(c.Context(), filter).Decode(&updateUser)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid email or password",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error finding user in database",
        })
    }

	if !utils.CheckPasswordHash(user.Password, updateUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid email or password",
        })
	}

	token, err := utils.GenerateJWT(user.Email)
	if err!= nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error generating JWT token",
        })
    }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token": token,
	})
}

func Logout(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	log.Println("logout token:", tokenStr)
	_,isExist := utils.IsTokenExpired(tokenStr)
	if!isExist {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Token not found or expired",
        })
    }
	utils.ClearExpiredTokens(tokenStr)
	c.ClearCookie()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Logout successful",
    })
}
