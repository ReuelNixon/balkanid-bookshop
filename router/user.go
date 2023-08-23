package router

import (
	db "bookshop/database"
	"bookshop/models"
	"bookshop/util"
	"math/rand"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(os.Getenv("PRIV_KEY"))

func SetupUserRoutes() {
    USER.Post("/signup", CreateUser)
}

func CreateUser(c *fiber.Ctx) error {
    u := new(models.User)

    if err := c.BodyParser(u); err != nil {
        return c.JSON(fiber.Map{
            "error": true,
            "input": "Please review your input",
        })
    }

    errors := util.ValidateRegister(u)
    if errors.Err {
        return c.JSON(errors)
    }

    if count := db.DB.Where(&models.User{Email: u.Email}).First(new(models.User)).RowsAffected; count > 0 {
        errors.Err, errors.Email = true, "Email is already registered"
    }
    if count := db.DB.Where(&models.User{Username: u.Username}).First(new(models.User)).RowsAffected; count > 0 {
        errors.Err, errors.Username = true, "Username is already registered"
    }
    if errors.Err {
        return c.JSON(errors)
    }

    password := []byte(u.Password)
    hashedPassword, err := bcrypt.GenerateFromPassword(
        password,
        rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost,
    )

    if err != nil {
        panic(err)
    }
    u.Password = string(hashedPassword)

    if err := db.DB.Create(&u).Error; err != nil {
        return c.JSON(fiber.Map{
            "error":   true,
            "general": "Something went wrong, please try again later. 😕",
        })
    }

    accessToken, refreshToken := util.GenerateTokens(u.UUID.String())
    accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
    c.Cookie(accessCookie)
    c.Cookie(refreshCookie)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}