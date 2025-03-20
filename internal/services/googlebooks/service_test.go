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

		res, err := s.GoogleBooks.SearchBooks(ctx, "harry potter", 10)
		require.NoError(t, err)

		assert.NotNil(t, res)
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
