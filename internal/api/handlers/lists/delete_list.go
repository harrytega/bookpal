package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/httperrors"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/labstack/echo/v4"
)

func DeleteListRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.DELETE("/:list_id", handler.DeleteList())
}

func (h *Handler) DeleteList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var body types.List
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		err := h.service.DeleteList(ctx, body.ListID.String(), body.UserID.String())

		if err != nil {
			return httperrors.ErrInternalServerDeletingList
		}

		return c.JSON(http.StatusNoContent, map[string]string{
			"message": "list has been deleted",
		})
	}
}
