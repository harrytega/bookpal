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
		userID := test.Fixtures().User1.ID

		res, err := s.Lists.CreateList(ctx, userID, "Favorites")
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetAllUserLists(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		userID := test.Fixtures().User1.ID

		res, err := s.Lists.GetAllUserLists(ctx, userID)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetListByID(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		listID := test.Fixtures().List1.ListID
		userID := test.Fixtures().User1.ID

		res, err := s.Lists.GetListByID(ctx, listID, userID)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestUpdateListName(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		listID := test.Fixtures().List1.ListID
		userID := test.Fixtures().User1.ID
		newListName := "the worst"

		err := s.Lists.UpdateListName(ctx, listID, userID, newListName)
		require.NoError(t, err)
	})
}

func TestDeleteList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		listID := test.Fixtures().List1.ListID
		userID := test.Fixtures().User1.ID

		err := s.Lists.DeleteList(ctx, listID, userID)
		require.NoError(t, err)
	})
}

func TestAddBookToList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		listID := test.Fixtures().List1.ListID
		userID := test.Fixtures().User1.ID
		bookID := test.Fixtures().Book1.BookID

		initialList, err := s.Lists.GetListByID(ctx, listID, userID)
		require.NoError(t, err)

		isInitiallyInList := false
		if initialList.R != nil && initialList.R.Books != nil {
			for _, book := range initialList.R.Books {
				if book.BookID == bookID {
					isInitiallyInList = true
					break
				}
			}
		}

		if isInitiallyInList {
			err := s.Lists.RemoveBookFromList(ctx, listID, userID, bookID)
			require.NoError(t, err)
		}

		err = s.Lists.AddBookToList(ctx, listID, userID, bookID)
		require.NoError(t, err)

		updatedList, err := s.Lists.GetListByID(ctx, listID, userID)
		require.NoError(t, err)

		require.NotNil(t, updatedList.R)
		require.NotNil(t, updatedList.R.Books)

		bookFound := false
		for _, book := range updatedList.R.Books {
			if book.BookID == bookID {
				bookFound = true
			}
		}
		assert.True(t, bookFound)
	})
}

func TestRemoveBookFromList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		listID := test.Fixtures().List1.ListID
		userID := test.Fixtures().User1.ID
		bookID := test.Fixtures().Book1.BookID

		errAddBook := s.Lists.AddBookToList(ctx, listID, userID, bookID)
		require.NoError(t, errAddBook)

		err := s.Lists.RemoveBookFromList(ctx, listID, userID, bookID)
		require.NoError(t, err)

		booksAfterRemoval, errAfter := s.Lists.GetAllBooksFromList(ctx, userID, listID)
		require.NoError(t, errAfter)

		bookStillExits := false
		for _, b := range booksAfterRemoval {
			if b.BookID == bookID {
				bookStillExits = true
				break
			}
		}
		require.False(t, bookStillExits)
	})
}

func TestGetAllBooksFromList(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		listID := test.Fixtures().List2.ListID
		userID := test.Fixtures().User1.ID
		bookID := test.Fixtures().Book2.BookID

		list, err := s.Lists.GetListByID(ctx, listID, userID)
		require.NoError(t, err)

		book, err := s.Books.GetBookByID(ctx, bookID)
		require.NoError(t, err)
		// list.AddBooks(ctx, s.DB, false, book)
		errAdding := list.AddBooks(ctx, s.DB, false, book)
		require.NoError(t, errAdding, "Failed to add book to list")

		err = list.Reload(ctx, s.DB)
		require.NoError(t, err)

		books, err := s.Lists.GetAllBooksFromList(ctx, userID, listID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(books))
	})
}
