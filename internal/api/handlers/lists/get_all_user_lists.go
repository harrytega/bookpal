package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/api/httperrors"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

func GetAllUserListsRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.GET("", handler.GetAllUserLists())
}

func (h *Handler) GetAllUserLists() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID := auth.UserFromContext(ctx).ID

		res, err := h.service.GetAllUserLists(ctx, userID)
		if err != nil {
			return httperrors.ErrInternalServerFailedFetchAllLists
		}

		lists := []*types.List{}
		for _, list := range res {
			books := []*types.BookInMyDb{}
			if list.R != nil && list.R.Books != nil {
				for _, book := range list.R.Books {
					convertedBook := &types.BookInMyDb{
						BookID: strfmt.UUID4(book.BookID),
						Author: &book.Author,
						Title:  &book.Title,
					}
					if book.BookDescription.Valid {
						convertedBook.BookDescription = book.BookDescription.String
					}

					if book.Genre.Valid {
						convertedBook.Genre = book.Genre.String
					}

					if book.Pages.Valid {
						convertedBook.Pages = SafeInt32(book.Pages.Int)
					}

					if book.Publisher.Valid {
						convertedBook.Publisher = book.Publisher.String
					}

					if book.Rating.Valid {
						convertedBook.Rating = SafeInt32(book.Rating.Int)
					}

					if book.UserNotes.Valid {
						convertedBook.UserNotes = book.UserNotes.String
					}
					books = append(books, convertedBook)
				}
			}

			convertedList := &types.List{
				Name:   &list.Name,
				ListID: (*strfmt.UUID)(&list.ListID),
				Books:  books,
				UserID: (*strfmt.UUID)(&userID),
			}
			lists = append(lists, convertedList)
		}

		response := &types.GetAllListsResponse{
			Data: lists,
		}
		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
