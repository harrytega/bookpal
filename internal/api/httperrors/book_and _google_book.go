package httperrors

import "net/http"

var (
	ErrBadRequestMissingGoogleBookID              = NewHTTPError(http.StatusBadRequest, "MISSING_GOOGLE_BOOK_ID", "Book ID is required")
	ErrBadRequestMissingSearchQuery               = NewHTTPError(http.StatusBadRequest, "MISSING_SEARCH_QUERY", "Search query is required")
	ErrBadRequestInvalidPageNumber                = NewHTTPError(http.StatusBadRequest, "INVALID_PAGE_NUMBER", "Page must be a positive number")
	ErrBadRequestInvalidPageSizeNumber            = NewHTTPError(http.StatusBadRequest, "INVALID_PAGE_SIZE_NUMBER", "Page size must be a number between 1 and 30")
	ErrNotFoundGoogleBookNotFound                 = NewHTTPError(http.StatusNotFound, "GOOGLE_BOOK_NOT_FOUND", "Book not found")
	ErrInternalServerFailedFetchGoogleBookDetails = NewHTTPError(http.StatusInternalServerError, "ERROR_GETTING_GOOGLE_BOOK_DETAILS", "Failed to get book details")
	ErrInternalServerFailedBookSearch             = NewHTTPError(http.StatusInternalServerError, "ERROR_SEARCHING_BOOK", "Failed to search books")

	ErrBadRequestMissingBookID              = NewHTTPError(http.StatusBadRequest, "MISSING_BOOK_ID", "Book ID is required")
	ErrNotFoundBookNotFound                 = NewHTTPError(http.StatusNotFound, "BOOK_NOT_FOUND", "Book not found")
	ErrInternalServerFailedFetchBookDetails = NewHTTPError(http.StatusInternalServerError, "ERROR_GETTING_BOOK_DETAILS", "Failed to get book details")
	ErrInternalServerFailedFetchUserBooks   = NewHTTPError(http.StatusInternalServerError, "ERROR_FETCHING_USER_BOOKS", "Failed to fetch user books")
	ErrInternalServerFailedAddGoogleBook    = NewHTTPError(http.StatusInternalServerError, "ERROR_ADDING_GOOGLE_BOOK", "Failed to add book")
	ErrInternalServerFailedRatingNotes      = NewHTTPError(http.StatusInternalServerError, "ERROR_ADDING_RATING_NOTES", "Failed to add rating/notes")
	ErrInternalServerFailedDeleteBook       = NewHTTPError(http.StatusInternalServerError, "ERROR_DELETING_USER_BOOK", "Failed to delete user book")

	ErrBadRequestInvalidRating = NewHTTPError(http.StatusBadRequest, "INVALID_RATING", "Rating mus be between 1 - 5")
)
