package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/api/httperrors"

	"github.com/labstack/echo/v4"
)

func DeleteListRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.DELETE("/:list_id", handler.DeleteList())
}

func (h *Handler) DeleteList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		listID := c.Param("list_id")
		userID := auth.UserFromContext(ctx).ID

		err := h.service.DeleteList(ctx, listID, userID)

		if err != nil {
			return httperrors.ErrInternalServerDeletingList
		}

		return c.JSON(http.StatusNoContent, map[string]string{
			"message": "list has been deleted",
		})
	}
}
