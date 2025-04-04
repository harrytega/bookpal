package books

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/labstack/echo/v4"
)

func PutUpdateBookRatingAndNotesRoute(s *api.Server) *echo.Route {
	handler := newHandler(s.Books)
	return s.Router.APIV1Book.PUT("/:book_id", handler.PutUpdateBookRatingAndNotesRoute())
}

func (h *Handler) PutUpdateBookRatingAndNotesRoute() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var body types.BookInMyDb
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}
		bookID := c.Param("book_id")
		if bookID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "book id is required",
			})
		}
		userID := auth.UserFromContext(ctx).ID
		rating := int(body.Rating)
		ratingPtr := &rating
		err := h.service.UpdateBookRatingAndNotes(ctx, bookID, userID, &body.UserNotes, ratingPtr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to add rating/notes" + err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Rating/Notes has been added to the book",
		})
	}
}
