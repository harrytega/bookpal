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

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/lists/"+list.ListID, nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusNoContent, res.Result().StatusCode)
	})
}
