package handlers

import (
	"net/http"
	"strings"
	"test-project/internal/api/auth"
	"test-project/internal/services/books"
	"test-project/internal/services/googlebooks"
	"test-project/internal/services/lists"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	booksService       *books.Service
	listsService       *lists.Service
	googleBooksService *googlebooks.Service
}

func NewHandler(booksService *books.Service, listsService *lists.Service, googleBooksService *googlebooks.Service) *Handler {
	return &Handler{
		booksService:       booksService,
		listsService:       listsService,
		googleBooksService: googleBooksService,
	}
}

func (h *Handler) RegisterRoutes(e echo.Echo) {
	api := e.Group("/api/v1")

	booksGroup := api.Group("/books")
	booksGroup.GET("/:id", h.GetBookByID)
	// booksGroup.GET("search", h.SearchBooksFromGoogleBooksAPI)
}

func (h *Handler) GetBookByID(c echo.Context) error {
	ctx := c.Request().Context()
	user := auth.UserFromContext(ctx)
	if user.ID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized access",
		})
	}
	bookID := c.Param("id")
	if bookID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Book ID is required",
		})
	}

	book, err := h.booksService.GetBookByID(ctx, bookID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Book not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch book: " + err.Error(),
		})
	}

	if book.UserID != user.ID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You don't have permission to this book",
		})
	}

	return c.JSON(http.StatusOK, book)
}

// func (h *Handler) SearchBooksFromGoogleBooksAPI(c echo.Context) error {
// 	ctx := c.Request().Context()
// 	query := c.QueryParam("q")
// 	if query == "" {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "Search query is required",
// 		})
// 	}

// }
