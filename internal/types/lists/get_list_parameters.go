// Code generated by go-swagger; DO NOT EDIT.

package lists

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetListParams creates a new GetListParams object
// no default values defined in spec.
func NewGetListParams() GetListParams {

	return GetListParams{}
}

// GetListParams contains all the bound params for the get list operation
// typically these are obtained from a http.Request
//
// swagger:parameters getList
type GetListParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*ID of list to return.
	  Required: true
	  In: path
	*/
	ListID strfmt.UUID `param:"list_id"`
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetListParams() beforehand.
func (o *GetListParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rListID, rhkListID, _ := route.Params.GetOK("list_id")
	if err := o.bindListID(rListID, rhkListID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetListParams) Validate(formats strfmt.Registry) error {
	var res []error

	// list_id
	// Required: true
	// Parameter is provided by construction from the route

	if err := o.validateListID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindListID binds and validates parameter ListID from path.
func (o *GetListParams) bindListID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("list_id", "path", "strfmt.UUID", raw)
	}
	o.ListID = *(value.(*strfmt.UUID))

	if err := o.validateListID(formats); err != nil {
		return err
	}

	return nil
}

// validateListID carries on validations for parameter ListID
func (o *GetListParams) validateListID(formats strfmt.Registry) error {

	if err := validate.FormatOf("list_id", "path", "uuid", o.ListID.String(), formats); err != nil {
		return err
	}
	return nil
}
