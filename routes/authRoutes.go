package routes

import (
	"example.com/user-service/controllers"
	"example.com/user-service/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router *fiber.App) {
	router.Post("/signup", middlewares.ValidateRequestBody([]string{"username", "email", "password"}), controllers.Signup)
	router.Post("/login", middlewares.ValidateRequestBody([]string{"email", "password"}), controllers.Login)
	router.Delete("/logout", controllers.Logout)

	router.Post("/addBill", middlewares.ValidateRequestBody([]string{"billTitle", "amount"}), controllers.AddBill)
	router.Get("/getBills", controllers.GetBills)
}
