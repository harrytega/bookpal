package googlebooks_test

import (
	"context"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchBooks(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, total, err := s.GoogleBooks.SearchBooks(ctx, "harry potter", 10, 1)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Greater(t, total, int64(0))

		if total > 10 {
			page2, page2Total, err := s.GoogleBooks.SearchBooks(ctx, "harry potter", 10, 2)
			require.NoError(t, err)
			assert.NotNil(t, page2)
			assert.Equal(t, total, page2Total)
		}
	})
}

func TestGetBookByID(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.GoogleBooks.GetBookByID(ctx, "XtekEncdTZcC")
		require.NoError(t, err)

		assert.NotNil(t, res)
	})
}
