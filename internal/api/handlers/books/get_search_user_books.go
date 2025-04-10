package books

import (
	"fmt"
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/api/httperrors"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

func GetSearchUserBooksRoute(s *api.Server) *echo.Route {
	handler := newHandler(s.Books)
	return s.Router.APIV1Book.GET("/search", handler.SearchUserBooks())
}

func (h *Handler) SearchUserBooks() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		query := c.QueryParam("query")
		if query == "" {
			return httperrors.ErrBadRequestMissingSearchQuery
		}
		userID := auth.UserFromContext(ctx).ID

		pageSize := 10
		page := 1
		pagination := new(types.Pagination)
		if err := c.Bind(pagination); err != nil {
			fmt.Print("failed")
		} else {
			if pagination.PageSize > 0 && pagination.PageSize <= 30 {
				pageSize = int(pagination.PageSize)
			}
			if pagination.CurrentPage > 0 {
				page = int(pagination.CurrentPage)
			}
		}

		if pagination.CurrentPage < 0 {
			pagination.CurrentPage = 1
		}
		if pagination.PageSize < 0 || pagination.PageSize > 30 {
			pagination.PageSize = 10
		}

		res, totalItems, err := h.service.SearchUserBooks(ctx, query, userID, pageSize, page)
		if err != nil {
			return httperrors.ErrInternalServerFailedBookSearch
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
			}
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

		response := &types.GetUserBooksResponse{
			Data:       convertedBooks,
			Pagination: pagination,
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
