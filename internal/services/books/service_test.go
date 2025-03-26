package books_test

import (
	"context"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBookByID(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.Books.GetBookByID(ctx, "f56eb34c-0ceb-401a-9f9d-c55402b2b3b9")
		require.NoError(t, err)

		assert.NotNil(t, res)
	})
}

func TestGetUserBooks(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.Books.GetUserBooks(ctx, "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestDeleteBook(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		err := s.Books.DeleteBook(ctx, "f56eb34c-0ceb-401a-9f9d-c55402b2b3b9", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		require.NoError(t, err)
	})
}

func TestUpdateBookRatingAndNotes(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		notes := "this is a test note"
		rating := 3
		err := s.Books.UpdateBookRatingAndNotes(ctx, "f56eb34c-0ceb-401a-9f9d-c55402b2b3b9", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5", &notes, &rating)
		require.NoError(t, err)
	})
}

func TestAddGoogleBook(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		err := s.Books.AddGoogleBook(ctx, "P3LFEAAAQBAJ", "76a79a2b-fbd8-45a0-b35b-671a28a87acf")
		require.NoError(t, err)
	})
}

func TestSearchUserBooks(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.Books.SearchUserBooks(ctx, "test", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetBooksByGenre(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.Books.GetBooksByGenre(ctx, "Romance", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetTopRatedBooks(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.Books.GetTopRatedBooks(ctx, 5, "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
