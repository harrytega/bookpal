package books_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTopRatedBooks(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token

		res := test.PerformRequest(t, s, "GET", "/api/v1/books/rated", nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
