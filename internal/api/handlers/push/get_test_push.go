package push

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"test-project/internal/api"
	"test-project/internal/api/auth"
	"test-project/internal/util"
)

func GetPushTestRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Push.GET("/test", getPushTestHandler(s))
}

func getPushTestHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		user := auth.UserFromEchoContext(c)

		err := s.Push.SendToUser(ctx, user, "Hello", "World")
		if err != nil {
			log.Debug().Err(err).Str("user_id", user.ID).Msg("Error while sending push to user.")
			return err
		}

		log.Debug().Str("user_id", user.ID).Msg("Successfully sent push message.")

		return c.String(http.StatusOK, "Success")
	}
}
