package httperrors

import "net/http"

var (
	ErrBadRequestMissingListID                = NewHTTPError(http.StatusBadRequest, "MISSING_LIST_ID", "List ID is required")
	ErrNotFoundListNotFound                   = NewHTTPError(http.StatusNotFound, "LIST_NOT_FOUND", "List nof found")
	ErrInternalServerDeletingBookFromList     = NewHTTPError(http.StatusInternalServerError, "ERROR_DELETING_BOOK_FROM_LIST", "Failed to remove book from list")
	ErrInternalServerFailedAddBookToList      = NewHTTPError(http.StatusInternalServerError, "ERROR_ADDING_BOOK_TO_LIST", "Failed to add book to list")
	ErrInternalServerDeletingList             = NewHTTPError(http.StatusInternalServerError, "ERROR_DELETING_LIST", "Failed to delete list")
	ErrInternalServerFailedFetchBooksFromList = NewHTTPError(http.StatusInternalServerError, "ERROR_FETCHING_BOOKS_FROM_LIST", "Failed to fetch books from list")
	ErrInternalServerFailedFetchList          = NewHTTPError(http.StatusInternalServerError, "ERROR_FETCHING_LIST", "Failed to fetch list")
	ErrInternalServerFailedCreateList         = NewHTTPError(http.StatusInternalServerError, "ERROR_CREATING_LIST", "Failed to create list")
	ErrInternalServerFailedUpdateListName     = NewHTTPError(http.StatusInternalServerError, "FAILED_TO_UPDATE_LIST_NAME", "Failed to update list name")
	ErrInternalServerFailedFetchAllLists      = NewHTTPError(http.StatusInternalServerError, "ERROR_FETCHING_ALL_USER_LISTS", "Failed to fetch all lists from user")
)
