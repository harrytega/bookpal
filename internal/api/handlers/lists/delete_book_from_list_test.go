package lists_test

import (
	"context"
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveBookFromList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		token := test.Fixtures().User1AccessToken1.Token
		list := test.Fixtures().List1
		book := test.Fixtures().Book1
		payload := test.GenericPayload{
			"list_id": list.ListID,
			"name":    list.Name,
			"user_id": list.UserID,
		}
		err := s.Lists.AddBookToList(ctx, list.ListID, list.UserID, book.BookID)
		require.NoError(t, err)

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/lists/"+list.ListID+"/"+book.BookID, payload, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusNoContent, res.Result().StatusCode)
	})
}
