package books

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"test-project/internal/models"
	"test-project/internal/services/googlebooks"
	"test-project/internal/util/db"

	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Service struct {
	db                 *sql.DB
	googleBooksService *googlebooks.Service
}

func NewService(db *sql.DB, googleBooksService *googlebooks.Service) *Service {
	return &Service{
		db:                 db,
		googleBooksService: googleBooksService,
	}
}

func (s *Service) GetBookByID(ctx context.Context, bookID string) (*models.Book, error) {
	logger := log.Ctx(ctx).With().Str("bookID", bookID).Logger()
	book, err := models.FindBook(ctx, s.db, bookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("book not found with ID: " + bookID)
			return nil, fmt.Errorf("Book not found with ID: %s", bookID)
		}
		logger.Error().Err(err).Msg("error fetching book")
		return nil, fmt.Errorf("error fetching book: %w", err)
	}

	logger.Info().
		Str("bookID", bookID).
		Msg("book succesfully fetched")

	return book, nil
}

func (s *Service) GetUserBooks(ctx context.Context, userID string) (models.BookSlice, error) {
	logger := log.Ctx(ctx).With().
		Str("userID", userID).
		Logger()

	books, err := models.Books(
		models.BookWhere.UserID.EQ(userID),
		qm.OrderBy("title ASC"),
	).All(ctx, s.db)

	if err != nil {
		logger.Error().Err(err).Msg("error fetching user books")
		return nil, fmt.Errorf("Error fetching user books: %w", err)
	}

	logger.Info().
		Int("bookCount", len(books)).
		Msg("Succesfully fetched all books from user")

	return books, nil
}

func (s *Service) DeleteBook(ctx context.Context, bookID, userID string) error {
	logger := log.Ctx(ctx).With().
		Str("bookID", bookID).
		Str("userID", userID).
		Logger()

	exists, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).Exists(ctx, s.db)

	if err != nil {
		logger.Error().Err(err).Msg("error checking book")
		return fmt.Errorf("error checking book: %w", err)
	}
	if !exists {
		logger.Warn().Msg("book not found or not owned by user")
		return fmt.Errorf("book not found or not owned by user")
	}

	rows, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).DeleteAll(ctx, s.db)

	if err != nil {
		logger.Error().Err(err).Msg("error deleting book")
		return fmt.Errorf("error deleting book: %w", err)
	}
	if rows == 0 {
		logger.Warn().Msg("no books were deleted")
		return fmt.Errorf("no books were deleted")
	}

	logger.Info().Msg("book succesfully deleted")
	return nil
}

func (s *Service) UpdateBookRatingAndNotes(ctx context.Context, bookID, userID string, userNotes *string, rating *int) error {
	loggerCtx := log.Ctx(ctx).With().
		Str("bookID", bookID).
		Str("userID", userID)

	if userNotes != nil {
		loggerCtx = loggerCtx.Str("userNotes", *userNotes)
	}

	if rating != nil {
		loggerCtx = loggerCtx.Int("rating", *rating)
	}

	logger := loggerCtx.Logger()

	book, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).One(ctx, s.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Msg("book not found or not owned by user")
			return fmt.Errorf("book not found or not owned by user")
		}
		logger.Error().Err(err).Msg("error fetching the book")
		return fmt.Errorf("error fetching the book: %w", err)
	}

	if userNotes != nil {
		book.UserNotes = null.StringFromPtr(userNotes)
	}

	if rating != nil {
		if *rating < 1 || *rating > 5 {
			logger.Warn().Msg("rating must be between 1 and 5")
			return fmt.Errorf("rating must be between 1 and 5")
		}
		book.Rating = null.IntFromPtr(rating)
	}

	_, err = book.Update(ctx, s.db, boil.Infer())
	if err != nil {
		logger.Error().Err(err).Msg("error updating the book")
		return fmt.Errorf("error updating the book: %w", err)
	}

	logComplete := logger.Info().Str("bookID", bookID)

	if userNotes != nil {
		logComplete = logComplete.Str("updatedNotes", *userNotes)
	}

	if rating != nil {
		logComplete = logComplete.Int("updatedRating", *rating)
	}

	logComplete.Msg("book updated successfully")
	return nil
}

func (s *Service) AddGoogleBook(ctx context.Context, googleID, userID string) error {
	logger := log.Ctx(ctx).With().
		Str("googleID", googleID).
		Str("userID", userID).
		Logger()

	googleBook, err := s.googleBooksService.GetBookByID(ctx, googleID)
	if err != nil {
		logger.Error().Err(err).Msg("error fetching book from the Google Books API")
		return fmt.Errorf("error fetching book from the Google Books API: %w", err)
	}

	if err := db.WithTransaction(ctx, s.db, func(tx boil.ContextExecutor) error {

		newBook := &models.Book{
			Title:           googleBook.BookDetails.Title,
			Author:          googleBook.BookDetails.Authors[0],
			Publisher:       null.StringFromPtr(&googleBook.BookDetails.Publisher),
			BookDescription: null.StringFrom(googleBook.BookDetails.Description),
			Genre:           null.StringFromPtr(&googleBook.BookDetails.Genre[0]),
			Pages:           null.IntFrom(googleBook.BookDetails.Pages),
			UserID:          userID,
		}

		if err := newBook.Insert(ctx, tx, boil.Infer()); err != nil {
			logger.Error().Err(err).Msg("error inserting book")
			return fmt.Errorf("error inserting book: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	logger.Info().
		Msg("book succesfully added to database")

	return nil
}

func (s *Service) SearchUserBooks(ctx context.Context, searchTerm, userID string) (models.BookSlice, error) {
	logger := log.Ctx(ctx).With().
		Str("searchTerm", searchTerm).
		Str("userID", userID).
		Logger()

	books, err := models.Books(
		models.BookWhere.UserID.EQ(userID),
		qm.Where("(title ILIKE ? OR author ILIKE ?)",
			"%"+searchTerm+"%", "%"+searchTerm+"%"),
		qm.OrderBy("title ASC"),
	).All(ctx, s.db)

	if err != nil {
		logger.Error().Err(err).Msg("error searching user books")
		return nil, fmt.Errorf("error searching user books: %w", err)
	}

	logger.Info().
		Int("bookCount", len(books)).
		Msg("books of user were succesfully fetched")

	return books, nil
}

func (s *Service) GetBooksByGenre(ctx context.Context, genre, userID string) (models.BookSlice, error) {
	logger := log.Ctx(ctx).With().
		Str("genre", genre).
		Str("userID", userID).
		Logger()

	books, err := models.Books(
		models.BookWhere.UserID.EQ(userID),
		qm.Where("genre ILIKE ?", genre),
		qm.OrderBy("title ASC"),
	).All(ctx, s.db)

	if err != nil {
		logger.Error().Err(err).Msg("error fetching books by genre")
		return nil, fmt.Errorf("error fetching books by genre: %w", err)
	}

	logger.Info().
		Int("bookCount", len(books)).
		Msg("books fetched succesfully")

	return books, nil
}

func (s *Service) GetTopRatedBooks(ctx context.Context, bookLimit int, userID string) (models.BookSlice, error) {
	logger := log.Ctx(ctx).With().
		Int("bookLimit", bookLimit).
		Str("userID", userID).
		Logger()

	if bookLimit < 0 {
		bookLimit = 10
	}

	books, err := models.Books(
		models.BookWhere.UserID.EQ(userID),
		models.BookWhere.Rating.IsNotNull(),
		qm.OrderBy("rating DESC"),
		qm.Limit(bookLimit),
	).All(ctx, s.db)

	if err != nil {
		logger.Error().Err(err).Msg("error fetching top rated books")
		return nil, fmt.Errorf("error fetching top rated books:%w", err)
	}

	logger.Info().
		Int("bookCount", len(books)).
		Msg("books fetched succesfully")

	return books, nil
}
