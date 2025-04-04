package books

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/labstack/echo/v4"
)

func PostAddGoogleBookRoute(s *api.Server) *echo.Route {
	handler := newHandler(s.Books)
	return s.Router.APIV1Book.POST("", handler.AddGoogleBook())
}

func (h *Handler) AddGoogleBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var body types.GoogleBook
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		userID := auth.UserFromContext(ctx).ID
		err := h.service.AddGoogleBook(ctx, *body.GoogleID, userID)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to add book" + err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"message": "Book has been added",
		})

	}
}
