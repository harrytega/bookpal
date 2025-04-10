package books

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/api/httperrors"

	"github.com/labstack/echo/v4"
)

func DeleteBookRoute(s *api.Server) *echo.Route {
	handler := newHandler(s.Books)
	return s.Router.APIV1Book.DELETE("/:book_id", handler.DeleteBook())
}

func (h *Handler) DeleteBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		bookID := c.Param("book_id")
		userID := auth.UserFromContext(ctx).ID
		if bookID == "" {
			return httperrors.ErrBadRequestMissingBookID
		}

		err := h.service.DeleteBook(ctx, bookID, userID)
		if err != nil {
			return httperrors.ErrInternalServerFailedDeleteBook
		}

		return c.JSON(http.StatusNoContent, "book has been deleted")
	}
}
