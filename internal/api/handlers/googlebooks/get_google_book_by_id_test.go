package googlebooks_test

import (
	"net/http"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetGoogleBookByID(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/api/v1/google/XtekEncdTZcC", nil, test.HeadersWithAuth(t, test.Fixtures().User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}
