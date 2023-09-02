package router

import (
	"bookshop/database"
	"bookshop/models"
	"bookshop/util"
	"strconv"
	"strings"

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
    privPaths.Post("/checkoutAll", CheckoutAll)
    
	BOOK.Get("/", GetPaginatedBooks)
    BOOK.Post("/searchTitle", SearchTitle)
    BOOK.Post("/searchAuthor", SearchAuthor)
	BOOK.Get("/:bookID/reviews", GetPaginatedReviews)
	BOOK.Get("/:bookID", GetBookByID)
    BOOK.Get("/:bookID/recommendations", GetRecommendations)
}

 
func GetPaginatedBooks(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	pageSize := c.Query("pageSize", "12")

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

    var books []models.Book
    for _, cartItem := range cartItems {
        var book models.Book
        if err := database.DB.Where("id = ?", cartItem.BookID).First(&book).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   true,
                "message": "Failed to retrieve book details",
            })
        }
        books = append(books, book)
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Cart items fetched successfully",
        "data":    books,
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

    isInCart := true
    // Check if the book is in the user's cart
    var cartItem models.Cart
    if err := database.DB.Where("book_id = ?", bookID).Where("user_id = ?", userID).First(&cartItem).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            isInCart = false
        } 
    }

    bookIDNum, err := strconv.Atoi(bookID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to parse book ID",
            "errors":   err,
        })
    }
    bookIDint := uint(bookIDNum)
    historyEntry := models.History{
        UserID: userID,
        BookID: bookIDint,
    }
    if err := database.DB.Create(&historyEntry).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to create history entry",
        })
    }

    // Remove the item from the cart
    if isInCart {
        if err := database.DB.Delete(&cartItem).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   true,
                "message": "Failed to remove item from cart",
                "isInCart": isInCart,
            })
        }
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Checkout successful",
    })
}

func CheckoutAll(c *fiber.Ctx) error {
    uuid := c.Locals("id").(string)
    userID, err := convertUUIDtoUserID(uuid)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to parse user ID",
        })
    }

    // Check if the book is in the user's cart
    var cartItems []models.Cart
    if err := database.DB.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error":   true,
                "message": "Cart is empty",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to retrieve cart items",
        })
    }

    if len(cartItems) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Cart is empty",
        })
    }

    // Create a history entry for each item in the cart
    for _, cartItem := range cartItems {
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
    }

    // Remove all items from the cart
    if err := database.DB.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to remove items from cart",
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

    var books []models.Book
    for _, historyEntry := range historyEntries {
        var book models.Book
        if err := database.DB.Where("id = ?", historyEntry.BookID).First(&book).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   true,
                "message": "Failed to retrieve book details",
            })
        }
        books = append(books, book)
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "History entries fetched successfully",
        "data":    books,
    })
}

func SearchTitle(c *fiber.Ctx) error {
    type TitleInput struct {
        Title string `json:"book_title"`
    } 

    var input TitleInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Invalid request body",
        })
    }

    if input.Title == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Missing book_title parameter",
        })
    }

    query := input.Title

    var books []models.Book
    if err := database.DB.Where("CONCAT(' ', LOWER(book_title)) ILIKE ?", "% "+strings.ToLower(query)+"%").
            Find(&books).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to fetch search results",
            "errors":   err,
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Search results fetched successfully",
        "data":    books,
    })
}

func SearchAuthor(c *fiber.Ctx) error {
    type AuthorInput struct {
        Author string `json:"book_author"`
    } 

    var input AuthorInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Invalid request body",
        })
    }

    if input.Author == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Missing book_author parameter",
        })
    }

    query := input.Author

    var books []models.Book
    if err := database.DB.Where("CONCAT(' ', LOWER(book_author)) ILIKE ?", "% "+strings.ToLower(query)+"%").
            Find(&books).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Failed to fetch search results",
            "errors":   err,
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "error":   false,
        "message": "Search results fetched successfully",
        "data":    books,
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
