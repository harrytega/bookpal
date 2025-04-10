package books_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteBook(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		book := test.Fixtures().Book2
		res := test.PerformRequest(t, s, "DELETE", "/api/v1/books/"+book.BookID, nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusNoContent, res.Result().StatusCode)
	})
}
