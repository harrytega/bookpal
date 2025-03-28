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
		bookID := test.Fixtures().Book1.BookID

		res, err := s.Books.GetBookByID(ctx, bookID)
		require.NoError(t, err)

		assert.NotNil(t, res)
	})
}

func TestGetUserBooks(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		userID := test.Fixtures().User1.ID

		res, err := s.Books.GetUserBooks(ctx, userID)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestDeleteBook(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		bookID := test.Fixtures().Book1.BookID
		userID := test.Fixtures().User1.ID

		err := s.Books.DeleteBook(ctx, bookID, userID)
		require.NoError(t, err)
	})
}

func TestUpdateBookRatingAndNotes(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		bookID := test.Fixtures().Book1.BookID
		userID := test.Fixtures().User1.ID
		notes := "this is a test note"
		rating := 3
		err := s.Books.UpdateBookRatingAndNotes(ctx, bookID, userID, &notes, &rating)
		require.NoError(t, err)
	})
}

func TestAddGoogleBook(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		googleID := test.Fixtures().GoogleBookSummary1.GoogleID
		userID := test.Fixtures().User1.ID

		err := s.Books.AddGoogleBook(ctx, googleID, userID)
		require.NoError(t, err)
	})
}

func TestSearchUserBooks(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		searchTerm := "test"
		userID := test.Fixtures().User1.ID

		res, err := s.Books.SearchUserBooks(ctx, searchTerm, userID)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetBooksByGenre(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		genre := "Romance"
		userID := test.Fixtures().User1.ID

		res, err := s.Books.GetBooksByGenre(ctx, genre, userID)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetTopRatedBooks(t *testing.T) {
	test.DotEnvLoadLocalOrSkipTest(t)

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		bookLimit := 5
		userID := test.Fixtures().User1.ID

		res, err := s.Books.GetTopRatedBooks(ctx, bookLimit, userID)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
