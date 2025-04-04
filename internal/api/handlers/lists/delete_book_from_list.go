package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/types"
	"test-project/internal/util"

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
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "list id is required",
			})
		}
		bookID := c.Param("book_id")
		if bookID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "book id is required",
			})
		}
		var body types.List
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		if err := h.service.RemoveBookFromList(ctx, listID, user.ID, bookID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to remove book from the list",
			})
		}

		return c.JSON(http.StatusNoContent, map[string]string{
			"message": "book was removed from list",
		})

	}
}
