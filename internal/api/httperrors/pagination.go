package httperrors

import "net/http"

var (
	ErrBadRequestInvalidPaginationParameters = NewHTTPError(http.StatusBadRequest, "INVALID_PAGINATION", "invalid pagination")
)
