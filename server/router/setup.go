package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

var USER fiber.Router
var ADMIN fiber.Router
var BOOK fiber.Router

var jwtKey = []byte(os.Getenv("PRIV_KEY"))

func SetupRoutes(app *fiber.App) {
	app.Get("/api/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to the BookShop, Built with 🖤 by ReuelNixon",
		})
	})
	
    api := app.Group("/api")

	BOOK = api.Group("/book")
	SetupBookRoutes()

    USER = api.Group("/user")
	SetupUserRoutes()

	ADMIN = api.Group("/admin")
	SetupAdminRoutes()
}