package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/labstack/echo/v4"
)

func PostAddBookToListRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.POST("/:list_id/books", handler.AddBookToList())
}

func (h *Handler) AddBookToList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		listID := c.Param("list_id")
		userID := auth.UserFromContext(ctx).ID
		var body types.BookInMyDb
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		if listID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "list ID required",
			})
		}

		if err := h.service.AddBookToList(ctx, listID, userID, string(body.BookID)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to add book to list" + err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"message": "Book has been added to the list",
		})
	}
}
