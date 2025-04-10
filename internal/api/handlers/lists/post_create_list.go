package lists

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/api/httperrors"
	"test-project/internal/services/lists"
	"test-project/internal/types"
	"test-project/internal/util"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *lists.Service
}

func NewHandler(service *lists.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func PostCreateListRoute(s *api.Server) *echo.Route {
	handler := NewHandler(s.Lists)
	return s.Router.APIV1Lists.POST("", handler.CreateList())
}

func (h *Handler) CreateList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID := auth.UserFromContext(ctx).ID
		var body types.ListRequest
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		res, err := h.service.CreateList(ctx, userID, *body.Name)
		if err != nil {
			return httperrors.ErrInternalServerFailedCreateList
		}

		response := &types.ListRequest{
			Name: &res.Name,
		}
		return util.ValidateAndReturn(c, http.StatusCreated, response)
	}
}
