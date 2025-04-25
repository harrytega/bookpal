package books

import (
	"net/http"
	"strconv"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/api/httperrors"
	"test-project/internal/services/books"
	"test-project/internal/util"

	"test-project/internal/types"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *books.Service
}

func newHandler(service *books.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func GetUserBooksRoute(s *api.Server) *echo.Route {
	handler := newHandler(s.Books)
	return s.Router.APIV1Book.GET("", handler.GetUserBooks())
}

func (h *Handler) GetUserBooks() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID := auth.UserFromContext(ctx).ID
		page := 1
		pageSize := 10

		if pageParam := c.QueryParam("page"); pageParam != "" {
			var err error
			page, err = strconv.Atoi(pageParam)
			if err != nil || page < 0 {
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

		res, totalItems, err := h.service.GetUserBooks(ctx, userID, pageSize, page)
		if err != nil {
			return httperrors.ErrInternalServerFailedFetchUserBooks
		}

		convertedBooks := []*types.BookInMyDb{}
		for _, book := range res {
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
				ImageLink:       book.ImageLink.String,
			}
			convertedBooks = append(convertedBooks, convertedBook)
		}

		totalPages := (totalItems + int64(pageSize) - 1) / int64(pageSize)

		pagination := &types.Pagination{
			CurrentPage:     int64(page),
			TotalPages:      totalPages,
			PageSize:        int64(pageSize),
			TotalItems:      totalItems,
			HasNextPage:     page < int(totalPages),
			HasPreviousPage: page > 1,
		}

		response := &types.GetUserBooksResponse{
			Data:       convertedBooks,
			Pagination: pagination,
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
