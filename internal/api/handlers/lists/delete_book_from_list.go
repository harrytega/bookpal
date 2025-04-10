package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/api/httperrors"

	"github.com/labstack/echo/v4"
)

func DeleteBookFromListRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.DELETE("/:list_id/:book_id", handler.RemoveBookFromList())
}

func (h *Handler) RemoveBookFromList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.UserFromContext(ctx)
		listID := c.Param("list_id")
		if listID == "" {
			return httperrors.ErrBadRequestMissingListID
		}
		bookID := c.Param("book_id")
		if bookID == "" {
			return httperrors.ErrBadRequestMissingBookID
		}

		if err := h.service.RemoveBookFromList(ctx, listID, user.ID, bookID); err != nil {
			return httperrors.ErrInternalServerDeletingBookFromList
		}

		return c.JSON(http.StatusNoContent, map[string]string{
			"message": "book was removed from list",
		})

	}
}
