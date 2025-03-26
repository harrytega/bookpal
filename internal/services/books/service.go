package books

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"test-project/internal/models"
	"test-project/internal/services/googlebooks"

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
	book, err := models.FindBook(ctx, s.db, bookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Book not found with ID: %s", bookID)
		}
		return nil, fmt.Errorf("error fetching book: %w", err)
	}
	return book, nil
}

func (s *Service) GetUserBooks(ctx context.Context, userID string) (models.BookSlice, error) {
	books, err := models.Books(
		models.BookWhere.UserID.EQ(userID),
		qm.OrderBy("title ASC"),
	).All(ctx, s.db)

	if err != nil {
		return nil, fmt.Errorf("Error fetching user books: %w", err)
	}

	return books, nil
}

func (s *Service) DeleteBook(ctx context.Context, bookID, userID string) error {
	exists, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).Exists(ctx, s.db)

	if err != nil {
		return fmt.Errorf("error checking book: %w", err)
	}
	if !exists {
		return fmt.Errorf("book not found or not owned by user")
	}

	rows, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).DeleteAll(ctx, s.db)

	if err != nil {
		return fmt.Errorf("error deleting book: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no books were deleted")
	}
	return nil
}

func (s *Service) UpdateBookRatingAndNotes(ctx context.Context, bookID, userID string, userNotes *string, rating *int) error {
	book, err := models.Books(
		models.BookWhere.BookID.EQ(bookID),
		models.BookWhere.UserID.EQ(userID),
	).One(ctx, s.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("book not found or not owned by user")
		}
		return fmt.Errorf("error fetching the book: %w", err)
	}

	if userNotes != nil {
		book.UserNotes = null.StringFromPtr(userNotes)
	}

	if rating != nil {
		if *rating < 1 || *rating > 5 {
			return fmt.Errorf("rating must be between 1 and 5")
		}
		book.Rating = null.IntFromPtr(rating)
	}

	_, err = book.Update(ctx, s.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("error updating the book: %w", err)
	}
	return nil
}

func (s *Service) AddGoogleBook(ctx context.Context, googleID, userID string) error {
	exists, err := models.Books(
		models.BookWhere.UserID.EQ(userID),
	).Exists(ctx, s.db)

	if err != nil {
		return fmt.Errorf("error checking book existance: %w", err)
	}
	if exists {
		return fmt.Errorf("book already exits in user's library")
	}

	googleBook, err := s.googleBooksService.GetBookByID(ctx, googleID)
	if err != nil {
		return fmt.Errorf("error fetching book from the Google Books API: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

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
		return fmt.Errorf("error inserting book: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}
	return nil
}

func (s *Service) SearchUserBooks(ctx context.Context, searchTerm, userID string) (models.BookSlice, error) {
	books, err := models.Books(
		models.BookWhere.UserID.EQ(userID),
		qm.Where("(title ILIKE ? OR author ILIKE ?)",
			"%"+searchTerm+"%", "%"+searchTerm+"%"),
		qm.OrderBy("title ASC"),
	).All(ctx, s.db)

	if err != nil {
		return nil, fmt.Errorf("error searching user books: %w", err)
	}
	return books, nil
}

func (s *Service) GetBooksByGenre(ctx context.Context, genre, userID string) (models.BookSlice, error) {
	books, err := models.Books(
		models.BookWhere.UserID.EQ(userID),
		qm.Where("genre ILIKE ?", genre),
		qm.OrderBy("title ASC"),
	).All(ctx, s.db)

	if err != nil {
		return nil, fmt.Errorf("error fetching books by genre: %w", err)
	}

	return books, nil
}

func (s *Service) GetTopRatedBooks(ctx context.Context, bookLimit int, userID string) (models.BookSlice, error) {
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
		return nil, fmt.Errorf("error fetching top rated books:%w", err)
	}

	return books, nil
}
