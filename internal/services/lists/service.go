package lists

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"test-project/internal/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) createList(ctx context.Context, userID, listName string) (*models.List, error) {

	if listName == "" {
		return nil, errors.New("list name cannot be empty")
	}

	exists, err := models.Lists(
		models.ListWhere.Name.EQ(listName),
		models.ListWhere.UserID.EQ(userID),
	).Exists(ctx, s.db)

	if err != nil {
		return nil, fmt.Errorf("error checking list existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("a list with this name already exits")
	}

	newList := &models.List{
		Name:   listName,
		UserID: userID,
	}

	if err := newList.Insert(ctx, s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("error creating list: %w", err)
	}

	return newList, nil
}

func (s *Service) GetAllUserLists(ctx context.Context, userID string) (models.ListSlice, error) {
	lists, err := models.Lists(
		models.ListWhere.UserID.EQ(userID),
		qm.OrderBy("name ASC"),
	).All(ctx, s.db)

	if err != nil {
		return nil, fmt.Errorf("error fetching user lists: %w", err)
	}

	return lists, nil
}

func (s *Service) GetListByID(ctx context.Context, listID, userID string) (*models.List, error) {
	list, err := models.Lists(
		models.ListWhere.ListID.EQ(listID),
		models.ListWhere.UserID.EQ(userID),
		qm.Load("Books"),
	).One(ctx, s.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("list not found or not owned by user")
		}
		return nil, fmt.Errorf("error fetching list: %w", err)
	}

	return list, nil
}

func (s *Service) UpdateListName(ctx context.Context, listID, userID, newListName string) error {
	if newListName == "" {
		return errors.New("list name cannot be empty")
	}

	list, err := s.GetListByID(ctx, listID, userID)

	if err != nil {
		return err
	}

	list.Name = newListName

	_, err = list.Update(ctx, s.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("error updating list: %w", err)
	}

	return nil
}

func (s *Service) DeleteList(ctx context.Context, listID, userID string) error {
	list, err := s.GetListByID(ctx, listID, userID)

	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	if err := list.RemoveBooks(ctx, tx, list.R.Books...); err != nil {
		return fmt.Errorf("error removing books from list: %w", err)
	}

	if _, err := list.Delete(ctx, tx); err != nil {
		return fmt.Errorf("error deleting list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	return nil
}

func (s *Service) AddBookToList(ctx context.Context, listID, userID, bookID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error startig transaction %w", err)
	}
	defer tx.Rollback()

	list, err := models.Lists(
		models.ListWhere.ListID.EQ(listID),
		models.ListWhere.UserID.EQ(userID),
		qm.Load("Books"),
	).One(ctx, tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("list not found or not owned by user")
		}
		return fmt.Errorf("error fetching list: %w", err)
	}

	book, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).One(ctx, tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("book not found or not owned by user")
		}
		return fmt.Errorf("error fetching book: %w", err)
	}

	for _, b := range list.R.Books {
		if b.BookID == bookID {
			return errors.New("book is already in this list")
		}
	}

	if err := list.AddBooks(ctx, tx, true, book); err != nil {
		return fmt.Errorf("error adding book to list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}
	return nil
}

func (s *Service) RemoveBookFromList(ctx context.Context, userID, bookID, listID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	list, err := s.GetListByID(ctx, listID, userID)
	if err != nil {
		return nil
	}

	book, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).One(ctx, tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("book not found or not owned by user")
		}
		return fmt.Errorf("error fetching book: %w", err)
	}

	found := false
	for _, b := range list.R.Books {
		if b.BookID == bookID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("book is not on the list")
	}

	if err := list.RemoveBooks(ctx, tx, book); err != nil {
		return fmt.Errorf("error removing book from list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}
	return nil
}

func (s *Service) GetAllBooksFromList(ctx context.Context, userID, listID string) (models.BookSlice, error) {
	list, err := s.GetListByID(ctx, listID, userID)
	if err != nil {
		return nil, err
	}

	return list.R.Books, nil
}
