package router

import (
	db "bookshop/database"
	"bookshop/models"
	"bookshop/util"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)
 
func SetupAdminRoutes() {
    ADMIN.Post("/signin", LoginAdmin)
    ADMIN.Get("/get-access-token", GetAdminAccessToken)

    privAdmin := ADMIN.Group("/private")
	privAdmin.Use(util.SecureAuth())
	privAdmin.Get("/data", GetAdminData)
	privAdmin.Post("/logout", LogoutAdmin)
	privAdmin.Delete("/delete", DeleteAdmin)
}


func LoginAdmin(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	// check if a user exists
	u := new(models.Admin)
	if res := db.DB.Where(
		&models.Admin{Email: input.Identity}).Or(
		&models.Admin{Username: input.Identity},
	).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := util.GenerateTokens(u.UUID.String())
	accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func GetAdminData(c *fiber.Ctx) error {
	id := c.Locals("id")
	u := new(models.Admin)
	if res := db.DB.Where("uuid = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Cannot find the User"})
	}

	return c.JSON(u)
}

func GetAdminAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	refreshClaims := new(models.Claims)
	token, _ := jwt.ParseWithClaims(refreshToken, refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if res := db.DB.Where(
		"expires_at = ? AND issued_at = ? AND issuer = ?",
		refreshClaims.ExpiresAt, refreshClaims.IssuedAt, refreshClaims.Issuer,
	).First(&models.Claims{}); res.RowsAffected <= 0 {
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			c.ClearCookie("access_token", "refresh_token")
			return c.SendStatus(fiber.StatusForbidden)
		}
	} else {
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	_, accessToken := util.GenerateAccessClaims(refreshClaims.Issuer)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"access_token": accessToken})
}

func LogoutAdmin(c *fiber.Ctx) error {
	c.ClearCookie("access_token", "refresh_token")
	return c.SendStatus(fiber.StatusOK)
}

func DeleteAdmin(c *fiber.Ctx) error {
	id := c.Locals("id")
	u := new(models.Admin)
	if res := db.DB.Where("uuid = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Cannot find the User"})
	}

	db.DB.Delete(&u)
	return c.SendStatus(fiber.StatusOK)
}