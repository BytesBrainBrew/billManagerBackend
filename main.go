package main

import (
	"fmt"

	"example.com/user-service/config" // Import the correct package path
	"example.com/user-service/models"
	"example.com/user-service/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	fmt.Println("Welcome to the Server")

	// Initialize MongoDB connection
	config.ConnectDB()

	models.CreateIndexModel()
	models.CreateBillIndexModel()

	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:8080, http://192.168.1.89, http://192.168.1.89:8080, http://192.168.1.89:3000",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Access-Control-Allow-Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	routes.AuthRoutes(app)

	app.Listen(":8080")
}
