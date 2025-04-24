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

		pageSize := 10
		page := 1

		pageSizeStr := c.QueryParam("pageSize")
		if pageSizeStr != "" {
			pageSizeInt, err := strconv.ParseInt(pageSizeStr, 10, 64)
			if err == nil && pageSizeInt > 0 {
				if pageSizeInt <= 30 {
					pageSize = int(pageSizeInt)
				} else {
					pageSize = 30
				}
			}
		}

		pageStr := c.QueryParam("page")
		if pageStr != "" {
			pageInt, err := strconv.ParseInt(pageStr, 10, 64)
			if err == nil && pageInt > 0 {
				page = int(pageInt)
			}
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

		pagination := &types.Pagination{
			CurrentPage: int64(page),
			PageSize:    int64(pageSize),
			TotalItems:  totalItems,
			TotalPages:  (totalItems + int64(pageSize) - 1) / int64(pageSize),
		}
		pagination.HasNextPage = pagination.CurrentPage < pagination.TotalPages
		pagination.HasPreviousPage = pagination.CurrentPage > 1

		response := &types.GetGoogleBooksResponse{
			Data:       convertedBooks,
			Pagination: pagination,
		}
		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
