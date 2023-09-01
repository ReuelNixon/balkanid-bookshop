package main

import (
	"log"

	"bookshop/database"
	"bookshop/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func CreateServer() *fiber.App {
    app := fiber.New()
    return app
}

func main() {
    database.ConnectToDB()

    app := CreateServer()
    app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:8080",
        AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
        AllowCredentials: true,
    }))
    app.Use(logger.New())

    router.SetupRoutes(app)
    
    // 404
    app.Use(func(c *fiber.Ctx) error {
        return c.SendStatus(404)
    })

    log.Fatal(app.Listen(":3000"))
}