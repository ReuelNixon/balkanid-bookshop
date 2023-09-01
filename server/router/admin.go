package router

import (
	db "bookshop/database"
	"bookshop/models"
	"bookshop/util"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)
 
func SetupAdminRoutes() {
	ADMIN.Post("/signup", CreateAdmin)
    ADMIN.Post("/signin", LoginAdmin)
    ADMIN.Get("/get-access-token", GetAdminAccessToken)

    privAdmin := ADMIN.Group("/private")
	privAdmin.Use(util.SecureAuth())
	privAdmin.Get("/data", GetAdminData)
	privAdmin.Post("/logout", LogoutAdmin)
	privAdmin.Delete("/delete", DeleteAdmin)
	privAdmin.Post("/addBook", AddBook)
    privAdmin.Delete("/:bookID", DeleteBook)
}

func CreateAdmin(c *fiber.Ctx) error {
    a := new(models.Admin)

    if err := c.BodyParser(a); err != nil {
        return c.JSON(fiber.Map{
            "error": true,
            "input": "Please review your input",
        })
    }

    errors := util.ValidateAdminRegister(a)
    if errors.Err {
        return c.JSON(errors)
    }

    if err := db.DB.Where(&models.Admin{Email: a.Email}).First(new(models.Admin)).Error; err == nil {
		errors.Err, errors.Email = true, "Email is already registered"
	}
	if err := db.DB.Where(&models.Admin{Username: a.Username}).First(new(models.Admin)).Error; err == nil {
		errors.Err, errors.Username = true, "Username is already registered"
	}
    if errors.Err {
        return c.JSON(errors)
    }

    password := []byte(a.Password)
    hashedPassword, err := bcrypt.GenerateFromPassword(
        password,
        rand.Intn(5),
    )

    if err != nil {
        panic(err)
    }
    a.Password = string(hashedPassword)

    if err := db.DB.Create(&a).Error; err != nil {
        return c.JSON(fiber.Map{
            "error":   true,
            "general": "Something went wrong, please try again later. 😕",
        })
    }

    accessToken, refreshToken := util.GenerateTokens(a.UUID.String())
    accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
    c.Cookie(accessCookie)
    c.Cookie(refreshCookie)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
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
		c.Cookie(util.ClearCookie("access_token"))
		c.Cookie(util.ClearCookie("refresh_token"))
		return c.SendStatus(fiber.StatusForbidden)
	}

	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			c.Cookie(util.ClearCookie("access_token"))
			c.Cookie(util.ClearCookie("refresh_token"))
			return c.SendStatus(fiber.StatusForbidden)
		}
	} else {
		c.Cookie(util.ClearCookie("access_token"))
		c.Cookie(util.ClearCookie("refresh_token"))
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
	c.Cookie(util.ClearCookie("access_token"))
	c.Cookie(util.ClearCookie("refresh_token"))
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

func AddBook(c *fiber.Ctx) error {
    input := new(models.Book)
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Invalid request body",
        })
    }

	if input.ISBN == "" || input.BookTitle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ISBN and Book Title are required",
		})
	}

	var lastBook models.Book
	if err := db.DB.Last(&lastBook).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			input.ID = 1
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to retrieve book",
			})
		}
	}

	input.ID = lastBook.ID + 1

    // Create a new book
    if err := db.DB.Create(&input).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to add book",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Book added successfully",
    })
}

func DeleteBook(c *fiber.Ctx) error {
    bookID := c.Params("bookID")

    // Check if the book exists
    var book models.Book
    if err := db.DB.Where("ID = ?", bookID).First(&book).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error":   true,
                "message": "Book not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to retrieve book",
        })
    }

    // Delete the book
    if err := db.DB.Delete(&book).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to delete book",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Book deleted successfully",
    })
}
