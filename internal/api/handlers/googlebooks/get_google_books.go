package googlebooks

import (
	"net/http"
	"strconv"
	"test-project/internal/api"
	"test-project/internal/services/googlebooks"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *googlebooks.Service
}

func NewHandler(service *googlebooks.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func GetGoogleBooksRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.GoogleBooks)
	return s.Router.APIV1Google.GET("/search", handler.SearchBooks())
}

func (h *Handler) SearchBooks() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		query := c.QueryParam("q")
		if query == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "search query is required",
			})
		}

		maxResults := 10
		if maxParam := c.QueryParam("maxResults"); maxParam != "" {
			var err error
			maxResults, err = strconv.Atoi(maxParam)
			if err != nil || maxResults < 1 || maxResults > 40 {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "maxResults must be a number between 1 and 40",
				})
			}
		}

		res, err := h.service.SearchBooks(ctx, query, maxResults)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to search books " + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, res)
	}
}
