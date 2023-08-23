package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

var USER fiber.Router
var ADMIN fiber.Router
var jwtKey = []byte(os.Getenv("PRIV_KEY"))

func SetupRoutes(app *fiber.App) {
    api := app.Group("/api")

    USER = api.Group("/user")
	SetupUserRoutes()

	ADMIN = api.Group("/admin")
	SetupAdminRoutes()
}