package lists_test

import (
	"context"
	"test-project/internal/api"
	"test-project/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateList(t *testing.T) {

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.Lists.CreateList(ctx, "f6ede5d8-e22a-4ca5-aa12-67821865a3e5", "Favorites")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetAllUserLists(t *testing.T) {

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.Lists.GetAllUserLists(ctx, "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetListByID(t *testing.T) {

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		res, err := s.Lists.GetListByID(ctx, "b9238b91-97e2-4837-97c5-a560761ffa81", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestUpdateListName(t *testing.T) {

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		err := s.Lists.UpdateListName(ctx, "b9238b91-97e2-4837-97c5-a560761ffa81", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5", "the worst")
		require.NoError(t, err)
	})
}

func TestDeleteList(t *testing.T) {

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		err := s.Lists.DeleteList(ctx, "b9238b91-97e2-4837-97c5-a560761ffa81", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		require.NoError(t, err)
	})
}

func TestAddBookToList(t *testing.T) {

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		err := s.Lists.AddBookToList(ctx, "b9238b91-97e2-4837-97c5-a560761ffa81", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5", "f56eb34c-0ceb-401a-9f9d-c55402b2b3b9")
		require.NoError(t, err)
	})
}

func TestRemoveBookFromList(t *testing.T) {

	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		err := s.Lists.RemoveBookFromList(ctx, "f6ede5d8-e22a-4ca5-aa12-67821865a3e5", "9238b91-97e2-4837-97c5-a560761ffa81", "f56eb34c-0ceb-401a-9f9d-c55402b2b3b9")
		require.NoError(t, err)
	})
}

func TestGetAllBooksFromList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		list, _ := s.Lists.GetListByID(ctx, "b9238b91-97e2-4837-97c5-a560761ffa81", "f6ede5d8-e22a-4ca5-aa12-67821865a3e5")
		book, _ := s.Books.GetBookByID(ctx, "f56eb34c-0ceb-401a-9f9d-c55402b2b3b9")
		list.AddBooks(ctx, s.DB, false, book)
		res, err := s.Lists.GetAllBooksFromList(ctx, "f6ede5d8-e22a-4ca5-aa12-67821865a3e5", list.ListID)

		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
