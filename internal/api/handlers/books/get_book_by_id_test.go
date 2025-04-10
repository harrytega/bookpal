package books_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBookByIDSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		bookID := test.Fixtures().Book1.BookID

		res := test.PerformRequest(t, s, "GET", "/api/v1/books/"+bookID, nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}

func TestGetBookByIDNotFound(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		nonExistentBookID := "71474c87-9ad0-4600-8d54-5ce10c19f2a4"

		resNotFound := test.PerformRequest(t, s, "GET", "/api/v1/books/"+nonExistentBookID, nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusNotFound, resNotFound.Result().StatusCode)
	})

}

func TestGetBookByIDInvalid(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		invalidBookID := "234lkdf,caksdaadsf99332"

		resInvalid := test.PerformRequest(t, s, "GET", "/api/v1/books/"+invalidBookID, nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusInternalServerError, resInvalid.Result().StatusCode)
	})
}
