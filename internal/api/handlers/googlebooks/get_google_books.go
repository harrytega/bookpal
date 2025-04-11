package googlebooks

import (
	"net/http"
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

		pageSize := 10
		page := 1
		pagination := new(types.Pagination)
		if err := c.Bind(pagination); err != nil {
			return httperrors.ErrBadRequestInvalidPaginationParameters
		}
		if pagination.PageSize > 0 && pagination.PageSize <= 30 {
			pageSize = int(pagination.PageSize)
		}
		if pagination.CurrentPage > 0 {
			page = int(pagination.CurrentPage)
		}

		res, totalItems, err := h.service.SearchBooks(ctx, query, pageSize, page)
		if err != nil {
			return httperrors.ErrInternalServerFailedBookSearch
		}

		convertedBooks := []*types.GoogleBook{}
		for _, book := range res.Books {
			convertedBook := &types.GoogleBook{
				GoogleBookID: &book.GoogleID,
				Title:        &book.BookDetails.Title,
			}
			if len(book.BookDetails.Authors) > 0 {
				convertedBook.Author = &book.BookDetails.Authors[0]
			} else {
				unknownAuthor := "Unknown"
				convertedBook.Author = &unknownAuthor
			}

			convertedBook.Publisher = book.BookDetails.Publisher
			convertedBook.BookDescription = book.BookDetails.Description

			if len(book.BookDetails.Genre) > 0 {
				convertedBook.Genre = book.BookDetails.Genre[0]
			}

			convertedBook.Pages = SafeInt32(book.BookDetails.Pages)
			convertedBooks = append(convertedBooks, convertedBook)
		}

		if pageSize <= 0 {
			pageSize = 1
		}

		pagination.CurrentPage = int64(page)
		pagination.PageSize = int64(pageSize)
		pagination.TotalPages = (totalItems + pagination.PageSize - 1) / pagination.PageSize
		pagination.HasNextPage = pagination.CurrentPage < pagination.TotalPages
		pagination.HasPreviousPage = pagination.CurrentPage > 1
		pagination.TotalItems = totalItems

		response := &types.GetGoogleBooksResponse{
			Data:       convertedBooks,
			Pagination: pagination,
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
