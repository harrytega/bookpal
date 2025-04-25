package googlebooks

import (
	"math"
	"net/http"
	"strings"
	"test-project/internal/api"
	"test-project/internal/api/httperrors"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/labstack/echo/v4"
)

func GetGoogleBookByIDRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.GoogleBooks)
	return s.Router.APIV1Google.GET("/:google_book_id", handler.GetBookByID())
}

func (h *Handler) GetBookByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("google_book_id")
		if id == "" {
			return httperrors.ErrBadRequestMissingGoogleBookID
		}

		book, err := h.service.GetBookByID(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return httperrors.ErrNotFoundGoogleBookNotFound
			}
			return httperrors.ErrInternalServerFailedFetchGoogleBookDetails
		}

		response := &types.GoogleBook{
			GoogleBookID:    &id,
			Title:           &book.BookDetails.Title,
			Author:          &book.BookDetails.Authors[0],
			Publisher:       book.BookDetails.Publisher,
			BookDescription: book.BookDetails.Description,
			Genre:           book.BookDetails.Genre[0],
			Pages:           SafeInt32(book.BookDetails.Pages),
			ImageLink:       book.BookDetails.ImageLinks.Thumbnail,
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
