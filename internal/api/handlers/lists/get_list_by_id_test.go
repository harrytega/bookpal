package lists_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetListByID(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		token := test.Fixtures().User1AccessToken1.Token
		listID := test.Fixtures().List1.ListID
		res := test.PerformRequest(t, s, "GET", "/api/v1/lists/"+listID, nil, test.HeadersWithAuth(t, token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
