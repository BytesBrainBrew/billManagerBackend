package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ValidateRequestBody(requiredFields []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
        body := make(map[string]interface{})
        if err := c.BodyParser(&body); err!= nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "message": "Error parsing request body",
            })
        }

        for _, field := range requiredFields {
            if _, ok := body[field];!ok {
                return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                    "message": fmt.Sprintf("Missing required field: %s", field),
                })
            }
        }

        c.Next()
        return nil
    }
}