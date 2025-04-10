package books_test

import (
	"encoding/json"
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchUserBooks(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		res := test.PerformRequest(t, s, "GET", "/api/v1/books/search?query=Test", nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

	})
}

func TestSearchUserBooksEmptyQuery(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		resEmptyQuery := test.PerformRequest(t, s, "GET", "/api/v1/books/search?query=", nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusBadRequest, resEmptyQuery.Result().StatusCode)
	})
}

func TestSearchUserBooksNonExistent(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		nonExistentBook := "Non"

		resNonExistentBook := test.PerformRequest(t, s, "GET", "/api/v1/books/search?query="+nonExistentBook, nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusOK, resNonExistentBook.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resNonExistentBook.Result().Body).Decode(&response)
		require.NoError(t, err, "Response should be valid JSON")

		data, ok := response["data"].([]interface{})
		require.True(t, ok, "Data should be an array")
		assert.Empty(t, data, "Search should return empty results for non-existent query")
	})
}
