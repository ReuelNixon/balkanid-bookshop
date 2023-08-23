package router

import (
	"github.com/gofiber/fiber/v2"
)

var USER fiber.Router

func SetupRoutes(app *fiber.App) {
    api := app.Group("/api")
    USER = api.Group("/user")
	SetupUserRoutes()
}