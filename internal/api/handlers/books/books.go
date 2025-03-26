package handlers

import (
	"net/http"
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
}

func (h *Handler) GetBookByID(c echo.Context) error {
	ctx := c.Request().Context()
	bookID := c.Param("id")

	book, err := h.booksService.GetBookByID(ctx, bookID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch book: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, book)
}
