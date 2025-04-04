package books_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBookByID(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		bookID := test.Fixtures().Book1.BookID
		res := test.PerformRequest(t, s, "GET", "/api/v1/books/"+bookID, nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
