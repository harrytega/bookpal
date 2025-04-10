package lists_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateListName(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		listID := test.Fixtures().List2.ListID
		userID := test.Fixtures().User1.ID
		payload := test.GenericPayload{
			"list_id": listID,
			"name":    "funny",
			"user_id": userID,
		}
		token := test.Fixtures().User1AccessToken1.Token
		res := test.PerformRequest(t, s, "PUT", "/api/v1/lists/"+listID, payload, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
