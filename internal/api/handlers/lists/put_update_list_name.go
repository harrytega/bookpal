package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/labstack/echo/v4"
)

func PutUpdateListNameRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.PUT("/:list_id", handler.UpdateListName())
}

func (h *Handler) UpdateListName() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID := auth.UserFromContext(ctx).ID
		listID := c.Param("list_id")
		if listID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "list id is required",
			})
		}
		var body types.List
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		err := h.service.UpdateListName(ctx, listID, userID, *body.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to update list name",
			})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": "List name has been changed",
		})
	}
}
