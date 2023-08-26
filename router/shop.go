package router

import (
	"bookshop/database"
	"bookshop/models"
	"bookshop/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupBookRoutes() {
	privPaths := BOOK.Group("/private")
	privPaths.Use(util.SecureAuth())
	privPaths.Post("/addToCart", AddToCart)
	privPaths.Get("/cart", GetCartItems)
	privPaths.Post("/postReview", PostReview)
    privPaths.Post("/:bookID/checkout", Checkout)
    privPaths.Get("/purchases", GetPurchases)
    
	BOOK.Get("/", GetPaginatedBooks)
	BOOK.Get("/:bookID/reviews", GetPaginatedReviews)
	BOOK.Get("/:bookID", GetBookByID)
}

 
func GetPaginatedBooks(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	pageSize := c.Query("pageSize", "10")

	var books []models.Book
	var totalCount int64

	// Count total number of books
	if err := database.DB.Model(&models.Book{}).Count(&totalCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve book count",
		})
	}

	// Calculate offset and limit based on pagination parameters
	offset, limit := calculatePagination(page, pageSize, totalCount)

	// Fetch paginated books
	if err := database.DB.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error":   false,
				"message": "No books found",
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve books",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Books fetched successfully",
		"data":    books,
	})
}


func GetBookByID(c *fiber.Ctx) error {
    bookID := c.Params("bookID")
    var book models.Book
    if err := database.DB.Where("ID = ?", bookID).First(&book).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error":   true,
                "message": "Book not found",
            })
        }
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Book details fetched successfully",
        "data":    book,
    })
}


func AddToCart(c *fiber.Ctx) error {
    input := new(models.Cart)

    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Invalid request body",
        })
    }
    if input.BookID == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Book ID missing in request body",
        })
    }

    if doesExist := database.DB.Where("id = ?", input.BookID).First(&models.Book{}).RowsAffected; doesExist == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Book does not exist",
        })
    }

    uuid := c.Locals("id").(string)
    userID, err := convertUUIDtoUserID(uuid)

    if inCart := database.DB.Where("book_id = ?", input.BookID).Where("user_id = ?", userID).First(&models.Cart{}).RowsAffected; inCart != 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Book already in cart",
        })
    }


    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to parse user ID",
        })
    }

    cartItem := models.Cart{
        UserID: userID,
        BookID: input.BookID,
    }
	if err := database.DB.Create(&cartItem).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to add book to cart",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Book added to cart successfully",
    })
}


func GetCartItems(c *fiber.Ctx) error {
    uuid := c.Locals("id").(string)
    userID, err := convertUUIDtoUserID(uuid)
    
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to parse user ID",
        })
    }

    noOfItems := database.DB.Model(&models.Cart{UserID: userID}).Find(&models.Cart{}).RowsAffected
    if noOfItems == 0 {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "error":   false,
            "message": "No items in cart",
            "data":    nil,
        })
    }
    var cartItems []models.Cart
    err = database.DB.Where("user_id = ?", userID).Find(&cartItems).Error
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to retrieve cart items",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Cart items fetched successfully",
        "data":    cartItems,
    })
}


func PostReview(c *fiber.Ctx) error {
    type ReviewInput struct {
        BookID uint   `json:"book_id"`
        Review string `json:"review"`
        Rating uint   `json:"rating"`
    } 

    var input ReviewInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Invalid request body",
        })
    }

    if input.BookID == 0 || input.Review == "" || input.Rating == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Missing review parameters",
        })
    }

    uuid := c.Locals("id").(string)
    userID, err := convertUUIDtoUserID(uuid)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to parse user ID",
        })
    }
    
    newReview := models.Review{
        UserID: userID,
        BookID: input.BookID,
        Review: input.Review,
        Rating: input.Rating,
    }
    if err := database.DB.Create(&newReview).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to post review",
        })
    } 

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Review posted successfully",
    })
}


func GetPaginatedReviews(c *fiber.Ctx) error {
    bookID := c.Params("bookID")
    page := c.Query("page", "1")
    pageSize := c.Query("pageSize", "10")

    pageNum := parsePageNumber(page)
    pageSizeNum := parsePageSize(pageSize)
    var reviews []models.Review
    if err := database.DB.Where("book_id = ?", bookID).
            Offset((pageNum - 1) * pageSizeNum).
            Limit(pageSizeNum).
            Find(&reviews).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to fetch reviews for the book",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Reviews fetched successfully",
        "data":    reviews,
    })
}


func Checkout(c *fiber.Ctx) error {
    bookID := c.Params("bookID")
    uuid := c.Locals("id").(string)
    userID, err := convertUUIDtoUserID(uuid)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to parse user ID",
        })
    }

    // Check if the book is in the user's cart
    var cartItem models.Cart
    if err := database.DB.Where("book_id = ?", bookID).Where("user_id = ?", userID).First(&cartItem).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error":   true,
                "message": "Book not found in cart",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to retrieve cart item",
        })
    }

    // Create a history entry
    historyEntry := models.History{
        UserID: userID,
        BookID: cartItem.BookID,
    }
    if err := database.DB.Create(&historyEntry).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to create history entry",
        })
    }

    // Remove the item from the cart
    if err := database.DB.Delete(&cartItem).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to remove item from cart",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Checkout successful",
    })
}


func GetPurchases(c *fiber.Ctx) error {
    uuid := c.Locals("id").(string)
    userID, err := convertUUIDtoUserID(uuid)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to parse user ID",
        })
    }

    var historyEntries []models.History
    if err := database.DB.Where("user_id = ?", userID).Find(&historyEntries).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to retrieve history entries",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "History entries fetched successfully",
        "data":    historyEntries,
    })
}






// ---------------------UTILITY FUNCTIONS---------------------
func calculatePagination(page, pageSize string, totalCount int64) (int, int) {
	pageNum := parsePageNumber(page)
	pageSizeNum := parsePageSize(pageSize)

	// Calculate offset and limit for pagination
	offset := (pageNum - 1) * pageSizeNum
	limit := pageSizeNum

	if offset >= int(totalCount) {
		offset = 0
	}

	return offset, limit
}

func parsePageNumber(page string) int {
	pageNum := 1
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}
	return pageNum
}

func parsePageSize(pageSize string) int {
	pageSizeNum := 10
	if size, err := strconv.Atoi(pageSize); err == nil && size > 0 {
		pageSizeNum = size
	}
	return pageSizeNum
}

func convertUUIDtoUserID(uuid string) (uint, error) {
    var user models.User
    if err := database.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
        return 0, err
    }
    return user.ID, nil
}
