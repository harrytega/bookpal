package lists

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"test-project/internal/models"

	"github.com/rs/zerolog/log"
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

func (s *Service) CreateList(ctx context.Context, userID, listName string) (*models.List, error) {
	logger := log.Ctx(ctx).With().
		Str("userID", userID).
		Str("listName", listName).
		Logger()

	if listName == "" {
		logger.Warn().Msg("list name cannot be empty")
		return nil, errors.New("list name cannot be empty")
	}

	exists, err := models.Lists(
		models.ListWhere.Name.EQ(listName),
		models.ListWhere.UserID.EQ(userID),
	).Exists(ctx, s.db)

	if err != nil {
		logger.Error().Err(err).Msg("error checking list existence")
		return nil, fmt.Errorf("error checking list existence: %w", err)
	}

	if exists {
		logger.Warn().Msg("a list with this name already exists")
		return nil, fmt.Errorf("a list with this name already exists")
	}

	newList := &models.List{
		Name:   listName,
		UserID: userID,
	}

	if err := newList.Insert(ctx, s.db, boil.Infer()); err != nil {
		logger.Error().Err(err).Msg("error creating list")
		return nil, fmt.Errorf("error creating list: %w", err)
	}

	logger.Info().
		Str("listID", newList.ListID).
		Msg("List created successfully")

	return newList, nil
}

func (s *Service) GetAllUserLists(ctx context.Context, userID string) (models.ListSlice, error) {
	logger := log.Ctx(ctx).With().
		Str("userID", userID).
		Logger()

	lists, err := models.Lists(
		models.ListWhere.UserID.EQ(userID),
		qm.OrderBy("name ASC"),
	).All(ctx, s.db)

	if err != nil {
		logger.Error().Err(err).Msg("error fetching user lists")
		return nil, fmt.Errorf("error fetching user lists: %w", err)
	}

	logger.Info().Msg("Succesfully fetched all lists from user")
	return lists, nil
}

func (s *Service) GetListByID(ctx context.Context, listID, userID string) (*models.List, error) {
	logger := log.Ctx(ctx).With().
		Str("listID", listID).
		Str("userID", userID).
		Logger()

	list, err := models.Lists(
		models.ListWhere.ListID.EQ(listID),
		models.ListWhere.UserID.EQ(userID),
		qm.Load("Books"),
	).One(ctx, s.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("list not found or not owned by user")
			return nil, fmt.Errorf("list not found or not owned by user")
		}
		logger.Error().Err(err).Msg("error fetching list")
		return nil, fmt.Errorf("error fetching list: %w", err)
	}

	logger.Info().
		Str("listID", listID).
		Msg("list succesfully fetched")
	return list, nil
}

func (s *Service) UpdateListName(ctx context.Context, listID, userID, newListName string) error {
	logger := log.Ctx(ctx).With().
		Str("listID", listID).
		Str("userID", userID).
		Str("newListName", newListName).
		Logger()

	if newListName == "" {
		logger.Warn().Msg("list name cannot be empty")
		return errors.New("list name cannot be empty")
	}

	list, err := s.GetListByID(ctx, listID, userID)

	if err != nil {
		return err
	}

	list.Name = newListName

	_, err = list.Update(ctx, s.db, boil.Infer())
	if err != nil {
		logger.Error().Err(err).Msg("error updating list")
		return fmt.Errorf("error updating list: %w", err)
	}

	logger.Info().
		Str("listID", listID).
		Msg("List updated succesfully")

	return nil
}

func (s *Service) DeleteList(ctx context.Context, listID, userID string) error {
	logger := log.Ctx(ctx).With().
		Str("listID", listID).
		Str("userID", userID).
		Logger()

	list, err := s.GetListByID(ctx, listID, userID)

	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error starting transaction")
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	if err := list.RemoveBooks(ctx, tx, list.R.Books...); err != nil {
		logger.Error().Err(err).Msg("error removing books from list")
		return fmt.Errorf("error removing books from list: %w", err)
	}

	if _, err := list.Delete(ctx, tx); err != nil {
		logger.Error().Err(err).Msg("error deleting list")
		return fmt.Errorf("error deleting list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("error commiting transaction")
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	logger.Info().Msg("list deleted succesfully")
	return nil
}

func (s *Service) AddBookToList(ctx context.Context, listID, userID, bookID string) error {
	logger := log.Ctx(ctx).With().
		Str("listID", listID).
		Str("userID", userID).
		Str("bookID", bookID).
		Logger()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error starting transaction")
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
			logger.Warn().Msg("list not found or not owned by user")
			return errors.New("list not found or not owned by user")
		}
		logger.Error().Err(err).Msg("error fetching list")
		return fmt.Errorf("error fetching list: %w", err)
	}

	book, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).One(ctx, tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("book not found or not owned by user")
			return errors.New("book not found or not owned by user")
		}
		logger.Error().Err(err).Msg("error fetching book")
		return fmt.Errorf("error fetching book: %w", err)
	}

	for _, b := range list.R.Books {
		if b.BookID == bookID {
			logger.Warn().Msg("book is already in this list")
			return errors.New("book is already in this list")
		}
	}

	if err := list.AddBooks(ctx, tx, false, book); err != nil {
		logger.Error().Err(err).Msg("error adding book to list")
		return fmt.Errorf("error adding book to list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("error commiting transaction")
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	logger.Info().
		Str("bookID", bookID).
		Str("listID", listID).
		Msg("book succesfully added to list")

	return nil
}

func (s *Service) RemoveBookFromList(ctx context.Context, userID, bookID, listID string) error {
	logger := log.Ctx(ctx).With().
		Str("listID", listID).
		Str("userID", userID).
		Str("bookID", bookID).
		Logger()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error starting transaction")
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	list, err := models.Lists(
		models.ListWhere.ListID.EQ(listID),
		models.ListWhere.UserID.EQ(userID),
		qm.Load("Books"),
	).One(ctx, tx)
	if err != nil {
		return err
	}

	book, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).One(ctx, tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("book not found or not owned by user")
			return errors.New("book not found or not owned by user")
		}
		logger.Error().Err(err).Msg("error fetching book")
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
		logger.Warn().Msg("book is not on the list")
		return errors.New("book is not on the list")
	}

	if err := list.RemoveBooks(ctx, tx, book); err != nil {
		logger.Error().Err(err).Msg("error removing book from list")
		return fmt.Errorf("error removing book from list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("error commiting transaction")
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	logger.Info().
		Str("bookID", bookID).
		Msg("book succesfully removed from list")

	return nil
}

func (s *Service) GetAllBooksFromList(ctx context.Context, userID, listID string) (models.BookSlice, error) {
	logger := log.Ctx(ctx).With().
		Str("userID", userID).
		Str("listID", listID).
		Logger()

	list, err := s.GetListByID(ctx, listID, userID)
	if err != nil {
		return nil, err
	}

	logger.Info().
		Int("bookCount", len(list.R.Books)).
		Msg("Succesfully fetched all books from list")

	return list.R.Books, nil
}
