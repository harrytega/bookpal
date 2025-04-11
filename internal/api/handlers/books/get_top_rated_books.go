package books

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/api/httperrors"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

func GetTopRatedBooksRoute(s *api.Server) *echo.Route {
	handler := newHandler(s.Books)
	return s.Router.APIV1Book.GET("/rated", handler.GetTopRatedBooks())
}

func (h *Handler) GetTopRatedBooks() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID := auth.UserFromContext(ctx).ID

		res, err := h.service.GetTopRatedBooks(ctx, userID)
		if err != nil {
			return httperrors.ErrInternalServerFailedFetchUserBooks
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
