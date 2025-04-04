package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

func GetAllBooksFromListRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.GET("/:list_id/books", handler.GetAllBooksFromList())
}

func (h *Handler) GetAllBooksFromList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID := auth.UserFromContext(ctx).ID
		listID := c.Param("list_id")

		if listID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "list id is required",
			})
		}

		res, err := h.service.GetAllBooksFromList(ctx, userID, listID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to fetch books from user" + err.Error(),
			})
		}

		convertedBooks := []*types.BookInMyDb{}

		for _, book := range res {
			convertedBook := &types.BookInMyDb{
				BookID:          strfmt.UUID4(book.BookID),
				Author:          &book.Author,
				BookDescription: book.BookDescription.String,
				Genre:           book.Genre.String,
				Pages:           SafeInt32(book.Pages.Int),
				Publisher:       book.Publisher.String,
				Rating:          SafeInt32(book.Rating.Int),
				Title:           &book.Title,
				UserNotes:       book.UserNotes.String,
			}
			convertedBooks = append(convertedBooks, convertedBook)
		}

		response := &types.GetUserBooksResponse{
			Data: convertedBooks,
		}
		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
