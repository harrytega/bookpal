package lists_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		list := test.Fixtures().List1
		token := test.Fixtures().User1AccessToken1.Token

		payload := test.GenericPayload{
			"list_id": list.ListID,
			"name":    list.Name,
			"user_id": list.UserID,
		}

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/lists/"+list.ListID, payload, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusNoContent, res.Result().StatusCode)
	})
}
