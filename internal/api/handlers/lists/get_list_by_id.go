package lists

import (
	"math"
	"net/http"
	"strings"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

func GetListByIDRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.GET("/:list_id", handler.GetListByID())
}

func (h *Handler) GetListByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		listID := c.Param("list_id")
		if listID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "list ID required",
			})
		}
		userID := auth.UserFromContext(ctx).ID
		res, err := h.service.GetListByID(ctx, listID, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]string{
					"error": "book not found",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to fetch book" + err.Error(),
			})
		}
		convertedBooks := []*types.BookInMyDb{}

		for _, book := range res.R.Books {
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

		response := &types.List{
			ListID: (*strfmt.UUID)(&res.ListID),
			Name:   &res.Name,
			UserID: (*strfmt.UUID)(&res.UserID),
			Books:  convertedBooks,
		}
		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}

func SafeInt32(val int) int32 {
	if val > math.MaxInt32 {
		return math.MaxInt32
	}
	if val < math.MinInt32 {
		return math.MinInt32
	}
	return int32(val)
}
