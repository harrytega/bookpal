package books_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddGoogleBook(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"googleID": "LipZDwAAQBAJ",
			"title":    "The Three Musketeers",
			"author":   "Alexandre Dumas",
		}
		token := test.Fixtures().User1AccessToken1.Token
		res := test.PerformRequest(t, s, "POST", "/api/v1/books", payload, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusCreated, res.Result().StatusCode)
	})
}
