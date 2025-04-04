package books_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateBookRatingAndNotes(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		bookToUpdate := test.Fixtures().Book1
		token := test.Fixtures().User1AccessToken1.Token
		payload := test.GenericPayload{
			"title":     "Test Title",
			"author":    "Testo",
			"rating":    4,
			"userNotes": "good",
		}
		res := test.PerformRequest(t, s, "PUT", "/api/v1/books/"+bookToUpdate.BookID, payload, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
