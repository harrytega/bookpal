package googlebooks

import (
	"net/http"
	"strconv"
	"test-project/internal/api"
	"test-project/internal/api/httperrors"
	"test-project/internal/services/googlebooks"
	"test-project/internal/types"
	"test-project/internal/util"

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
		query := c.QueryParam("query")
		if query == "" {
			return httperrors.ErrBadRequestMissingSearchQuery
		}

		page := 1
		pageSize := 10

		if pageParam := c.QueryParam("page"); pageParam != "" {
			var err error
			page, err = strconv.Atoi(pageParam)
			if err != nil || page < 1 {
				return httperrors.ErrBadRequestInvalidPageNumber
			}
		}

		if pageSizeParam := c.QueryParam("pageSize"); pageSizeParam != "" {
			var err error
			pageSize, err = strconv.Atoi(pageSizeParam)
			if err != nil || pageSize < 1 || pageSize > 30 {
				return httperrors.ErrBadRequestInvalidPageSizeNumber
			}
		}

		res, totalItems, err := h.service.SearchBooks(ctx, query, pageSize, page)
		if err != nil {
			return httperrors.ErrInternalServerFailedBookSearch
		}

		convertedBooks := []*types.GoogleBook{}
		for _, book := range res.Books {
			convertedBook := &types.GoogleBook{
				GoogleID:        &book.GoogleID,
				Title:           &book.BookDetails.Title,
				Author:          &book.BookDetails.Authors[0],
				Publisher:       book.BookDetails.Publisher,
				BookDescription: book.BookDetails.Description,
				Genre:           book.BookDetails.Genre[0],
				Pages:           SafeInt32(book.BookDetails.Pages),
			}
			convertedBooks = append(convertedBooks, convertedBook)
		}

		totalPages := (totalItems + int64(pageSize) - 1) / int64(pageSize)

		pagination := &types.Pagination{
			CurrentPage:     int64(page),
			TotalPages:      totalPages,
			PageSize:        int64(pageSize),
			TotalItems:      totalItems,
			HasNextPage:     int64(page) < totalPages,
			HasPreviousPage: page > 1,
		}

		response := &types.GetGoogleBooksResponse{
			Data:       convertedBooks,
			Pagination: pagination,
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
