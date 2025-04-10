package lists_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddBookToList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		listID := test.Fixtures().List1.ListID
		book := test.Fixtures().Book1
		payload := test.GenericPayload{
			"book_id": book.BookID,
			"title":   book.Title,
			"author":  book.Author,
		}
		res := test.PerformRequest(t, s, "POST", "/api/v1/lists/"+listID+"/books", payload, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusCreated, res.Result().StatusCode)
	})
}
