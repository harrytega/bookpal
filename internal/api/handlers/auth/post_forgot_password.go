package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"test-project/internal/api"
	"test-project/internal/models"
	"test-project/internal/types"
	"test-project/internal/util"
	"test-project/internal/util/db"
)

func PostForgotPasswordRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/forgot-password", postForgotPasswordHandler(s))
}

func postForgotPasswordHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var body types.PostForgotPasswordPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		// enforce lowercase usernames, trim whitespaces
		username := util.ToUsernameFormat(body.Username.String())

		log := util.LogFromContext(ctx).With().Str("username", username).Logger()

		user, err := models.Users(models.UserWhere.Username.EQ(null.StringFrom(username))).One(ctx, s.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Debug().Str("username", username).Err(err).Msg("User not found")
				return c.NoContent(http.StatusNoContent)
			}

			log.Debug().Str("username", username).Err(err).Msg("Failed to load user")
			return err
		}

		if !user.IsActive {
			log.Debug().Msg("User is deactivated, rejecting password reset")
			return c.NoContent(http.StatusNoContent)
		}

		if !user.Password.Valid {
			log.Debug().Msg("User is missing password, forbidding password reset")
			return c.NoContent(http.StatusNoContent)
		}

		if s.Config.Auth.PasswordResetTokenDebounceDuration > 0 {
			resetTokenInDebounceTimeExists, err := user.PasswordResetTokens(
				models.PasswordResetTokenWhere.CreatedAt.GT(s.Clock.Now().Add(-s.Config.Auth.PasswordResetTokenDebounceDuration)),
				models.PasswordResetTokenWhere.ValidUntil.GT(s.Clock.Now()),
			).Exists(ctx, s.DB)
			if err != nil {
				log.Error().Err(err).Msg("Failed to check for existing password reset token")
				return err
			}

			if resetTokenInDebounceTimeExists {
				log.Debug().Msg("Password reset token exists within debounce time, not sending new one")
				return c.NoContent(http.StatusNoContent)
			}
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			passwordResetToken, err := user.PasswordResetTokens(
				models.PasswordResetTokenWhere.CreatedAt.GT(s.Clock.Now().Add(-s.Config.Auth.PasswordResetTokenReuseDuration)),
				models.PasswordResetTokenWhere.ValidUntil.GT(s.Clock.Now()),
			).One(ctx, tx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Debug().Err(err).Msg("No valid password reset token exists, creating new one")

					passwordResetToken = &models.PasswordResetToken{
						UserID:     user.ID,
						ValidUntil: s.Clock.Now().Add(s.Config.Auth.PasswordResetTokenValidity),
					}

					if err := passwordResetToken.Insert(ctx, tx, boil.Infer()); err != nil {
						log.Debug().Err(err).Msg("Failed to insert password reset token")
						return err
					}
				} else {
					log.Debug().Err(err).Msg("Failed to check for existing valid password reset token")
					return err
				}
			}

			u, err := url.Parse(s.Config.Frontend.BaseURL)
			if err != nil {
				log.Error().Err(err).Msg("Failed to parse frontend base URL")
				return err
			}

			u.Path = path.Join(u.Path, s.Config.Frontend.PasswordResetEndpoint)

			q := u.Query()
			q.Set("token", passwordResetToken.Token)
			u.RawQuery = q.Encode()

			if err := s.Mailer.SendPasswordReset(ctx, user.Username.String, u.String()); err != nil {
				log.Debug().Err(err).Msg("Failed to send password reset email")
				return err
			}

			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to initiate password reset")
			return err
		}

		log.Debug().Msg("Successfully initiated password reset")

		return c.NoContent(http.StatusNoContent)
	}
}
