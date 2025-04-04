package lists_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"list_ID": "b9238b91-97e2-4837-97c5-a560761ffa81",
			"name":    "uhhhh",
			"user_ID": test.Fixtures().User1.ID,
		}
		token := test.Fixtures().User1AccessToken1.Token
		res := test.PerformRequest(t, s, "POST", "/api/v1/lists", payload, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusCreated, res.Result().StatusCode)
	})
}
