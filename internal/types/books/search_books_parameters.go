// Code generated by go-swagger; DO NOT EDIT.

package books

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewSearchBooksParams creates a new SearchBooksParams object
// with the default values initialized.
func NewSearchBooksParams() SearchBooksParams {

	var (
		// initialize parameters with default values

		pageDefault     = int64(1)
		pageSizeDefault = int64(10)
	)

	return SearchBooksParams{
		Page: &pageDefault,

		PageSize: &pageSizeDefault,
	}
}

// SearchBooksParams contains all the bound params for the search books operation
// typically these are obtained from a http.Request
//
// swagger:parameters searchBooks
type SearchBooksParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Page number
	  In: query
	  Default: 1
	*/
	Page *int64 `query:"page"`
	/*amount of books per page
	  Maximum: 30
	  In: query
	  Default: 10
	*/
	PageSize *int64 `query:"pageSize"`
	/*Search term for searching books. (title, author or publisher)
	  Required: true
	  In: query
	*/
	Query string `query:"query"`
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSearchBooksParams() beforehand.
func (o *SearchBooksParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qPage, qhkPage, _ := qs.GetOK("page")
	if err := o.bindPage(qPage, qhkPage, route.Formats); err != nil {
		res = append(res, err)
	}

	qPageSize, qhkPageSize, _ := qs.GetOK("pageSize")
	if err := o.bindPageSize(qPageSize, qhkPageSize, route.Formats); err != nil {
		res = append(res, err)
	}

	qQuery, qhkQuery, _ := qs.GetOK("query")
	if err := o.bindQuery(qQuery, qhkQuery, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SearchBooksParams) Validate(formats strfmt.Registry) error {
	var res []error

	// page
	// Required: false
	// AllowEmptyValue: false

	// pageSize
	// Required: false
	// AllowEmptyValue: false

	if err := o.validatePageSize(formats); err != nil {
		res = append(res, err)
	}

	// query
	// Required: true
	// AllowEmptyValue: false
	if err := validate.Required("query", "query", o.Query); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindPage binds and validates parameter Page from query.
func (o *SearchBooksParams) bindPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewSearchBooksParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("page", "query", "int64", raw)
	}
	o.Page = &value

	return nil
}

// bindPageSize binds and validates parameter PageSize from query.
func (o *SearchBooksParams) bindPageSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewSearchBooksParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("pageSize", "query", "int64", raw)
	}
	o.PageSize = &value

	if err := o.validatePageSize(formats); err != nil {
		return err
	}

	return nil
}

// validatePageSize carries on validations for parameter PageSize
func (o *SearchBooksParams) validatePageSize(formats strfmt.Registry) error {

	// Required: false
	if o.PageSize == nil {
		return nil
	}

	if err := validate.MaximumInt("pageSize", "query", *o.PageSize, 30, false); err != nil {
		return err
	}

	return nil
}

// bindQuery binds and validates parameter Query from query.
func (o *SearchBooksParams) bindQuery(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("query", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("query", "query", raw); err != nil {
		return err
	}

	o.Query = raw

	return nil
}
