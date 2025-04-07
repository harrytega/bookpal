package books

import (
	"math"
	"net/http"
	"strings"
	"test-project/internal/api"
	"test-project/internal/api/httperrors"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

func GetBookByIDRoute(s *api.Server) *echo.Route {
	handler := newHandler(s.Books)
	return s.Router.APIV1Book.GET("/:book_id", handler.GetBookByID())
}

func (h *Handler) GetBookByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		bookID := c.Param("book_id")
		if bookID == "" {
			return httperrors.ErrBadRequestMissingBookID
		}

		res, err := h.service.GetBookByID(ctx, bookID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return httperrors.ErrNotFoundBookNotFound
			}
			return httperrors.ErrInternalServerFailedFetchBookDetails
		}

		response := &types.BookInMyDb{
			BookID:          strfmt.UUID4(res.BookID),
			Author:          &res.Author,
			BookDescription: res.BookDescription.String,
			Genre:           res.Genre.String,
			Pages:           SafeInt32(res.Pages.Int),
			Publisher:       res.Publisher.String,
			Rating:          SafeInt32(res.Rating.Int),
			Title:           &res.Title,
			UserNotes:       res.UserNotes.String,
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
