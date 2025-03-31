package googlebooks

import (
	"net/http"
	"strings"
	"test-project/internal/api"

	"github.com/labstack/echo/v4"
)

func GetGoogleBookByIDRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.GoogleBooks)
	return s.Router.APIV1Google.GET("/:id", handler.GetBookByID())
}

func (h *Handler) GetBookByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("id")
		if id != "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "book ID is required",
			})
		}

		book, err := h.service.GetBookByID(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]string{
					"error": "book not found",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get book details " + err.Error(),
			})
		}
		return c.JSON(http.StatusOK, book)
	}
}
