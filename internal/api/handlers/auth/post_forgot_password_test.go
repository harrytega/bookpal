package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"test-project/internal/api"
	"test-project/internal/api/httperrors"
	"test-project/internal/config"
	"test-project/internal/models"
	"test-project/internal/test"
	"test-project/internal/types"
	"test-project/internal/util/db"
)

func TestPostForgotPasswordSuccess(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	config.Auth.PasswordResetTokenReuseDuration = 120 * time.Second
	config.Auth.PasswordResetTokenDebounceDuration = 60 * time.Second

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
		require.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		passwordResetToken, err := fixtures.User1.PasswordResetTokens().One(ctx, s.DB)
		require.NoError(t, err)

		mail := test.GetLastSentMail(t, s.Mailer)
		require.NotNil(t, mail)
		assert.Contains(t, string(mail.HTML), fmt.Sprintf("http://localhost:3000/set-new-password?token=%s", passwordResetToken.Token))

		// retrying should not send a new mail because of the debounce time
		{
			res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
			require.Equal(t, http.StatusNoContent, res.Result().StatusCode)

			sentMails := test.GetSentMails(t, s.Mailer)
			assert.Len(t, sentMails, 1)
		}

		// CreatedAt of token exceeds debounce time, retrying should send a new mail
		// but with the same token as the reuse duration has not passed yet
		{
			passwordResetToken.CreatedAt = s.Clock.Now().Add(-s.Config.Auth.PasswordResetTokenDebounceDuration)
			_, err = passwordResetToken.Update(ctx, s.DB, boil.Whitelist(models.PasswordResetTokenColumns.CreatedAt))
			require.NoError(t, err)

			res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
			require.Equal(t, http.StatusNoContent, res.Result().StatusCode)

			sentMails := test.GetSentMails(t, s.Mailer)
			require.Len(t, sentMails, 2)

			passwordResetTokens, err := fixtures.User1.PasswordResetTokens().All(ctx, s.DB)
			require.NoError(t, err)

			assert.Len(t, passwordResetTokens, 1)
			for _, mail := range sentMails {
				assert.Contains(t, string(mail.HTML), fmt.Sprintf("http://localhost:3000/set-new-password?token=%s", passwordResetTokens[0].Token))
			}
		}

		// CreatedAt of token exceeds reuse time, retrying should send a new mail with a new token
		{
			passwordResetToken.CreatedAt = s.Clock.Now().Add(-s.Config.Auth.PasswordResetTokenReuseDuration)
			_, err = passwordResetToken.Update(ctx, s.DB, boil.Whitelist(models.PasswordResetTokenColumns.CreatedAt))
			require.NoError(t, err)

			res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
			require.Equal(t, http.StatusNoContent, res.Result().StatusCode)

			sentMails := test.GetSentMails(t, s.Mailer)
			require.Len(t, sentMails, 3)

			passwordResetTokens, err := fixtures.User1.PasswordResetTokens(
				db.OrderBy(types.OrderDirDesc, models.PasswordResetTokenColumns.CreatedAt),
			).All(ctx, s.DB)
			require.NoError(t, err)

			require.Len(t, passwordResetTokens, 2)

			assert.Contains(t, string(sentMails[2].HTML), fmt.Sprintf("http://localhost:3000/set-new-password?token=%s", passwordResetTokens[0].Token))
		}

		// Token validity is expired, retrying should send a new mail with a new token
		{
			_, err = models.PasswordResetTokens().UpdateAll(ctx, s.DB, models.M{
				models.PasswordResetTokenColumns.ValidUntil: s.Clock.Now().Add(-1 * time.Second),
			})
			require.NoError(t, err)

			res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
			require.Equal(t, http.StatusNoContent, res.Result().StatusCode)

			sentMails := test.GetSentMails(t, s.Mailer)
			require.Len(t, sentMails, 4)

			passwordResetTokens, err := fixtures.User1.PasswordResetTokens(
				db.OrderBy(types.OrderDirDesc, models.PasswordResetTokenColumns.CreatedAt),
			).All(ctx, s.DB)
			require.NoError(t, err)

			require.Len(t, passwordResetTokens, 3)

			assert.Contains(t, string(sentMails[3].HTML), fmt.Sprintf("http://localhost:3000/set-new-password?token=%s", passwordResetTokens[0].Token))
		}
	})
}

func TestPostForgotPasswordSuccessLowercaseTrimWhitespaces(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fmt.Sprintf(" %s ", strings.ToUpper(fixtures.User1.Username.String)),
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		passwordResetToken, err := fixtures.User1.PasswordResetTokens().One(ctx, s.DB)
		require.NoError(t, err)

		mail := test.GetLastSentMail(t, s.Mailer)
		require.NotNil(t, mail)
		assert.Contains(t, string(mail.HTML), fmt.Sprintf("http://localhost:3000/set-new-password?token=%s", passwordResetToken.Token))
	})
}

func TestPostForgotPasswordUnknownUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		payload := test.GenericPayload{
			"username": "definitelydoesnotexist@example.com",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := test.GetLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordDeactivatedUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		payload := test.GenericPayload{
			"username": fixtures.UserDeactivated.Username,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := test.GetLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordUserWithoutPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		payload := test.GenericPayload{
			"username": fixtures.User2.Username,
		}

		fixtures.User2.Password = null.String{}
		rowsAff, err := fixtures.User2.Update(context.Background(), s.DB, boil.Infer())
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAff)

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := test.GetLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordBadRequest(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		tests := []struct {
			name    string
			payload test.GenericPayload
		}{
			{
				name:    "MissingUsername",
				payload: test.GenericPayload{},
			},
			{
				name: "EmptyUsername",
				payload: test.GenericPayload{
					"username": "",
				},
			},
			{
				name: "InvalidUsername",
				payload: test.GenericPayload{
					"username": "definitely not an email",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", tt.payload, nil)
				assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

				var response httperrors.HTTPValidationError
				test.ParseResponseAndValidate(t, res, &response)

				test.Snapshoter.Save(t, response)

				cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)

				mail := test.GetLastSentMail(t, s.Mailer)
				assert.Nil(t, mail)
			})
		}
	})
}
