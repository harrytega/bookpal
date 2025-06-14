// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// List list
//
// swagger:model list
type List struct {

	// books
	Books []*BookInMyDb `json:"books"`

	// list id
	// Example: d6764ee3-bf09-40c3-97c5-8f78b7de7ec3
	// Required: true
	// Format: uuid
	ListID *strfmt.UUID `json:"list_id"`

	// name
	// Example: Cookbooks
	// Required: true
	Name *string `json:"name"`

	// user id
	// Example: d6764ee3-bf09-40c3-97c5-8f78b7de7ec3
	// Required: true
	// Format: uuid
	UserID *strfmt.UUID `json:"user_id"`
}

// Validate validates this list
func (m *List) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBooks(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateListID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *List) validateBooks(formats strfmt.Registry) error {
	if swag.IsZero(m.Books) { // not required
		return nil
	}

	for i := 0; i < len(m.Books); i++ {
		if swag.IsZero(m.Books[i]) { // not required
			continue
		}

		if m.Books[i] != nil {
			if err := m.Books[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("books" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("books" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *List) validateListID(formats strfmt.Registry) error {

	if err := validate.Required("list_id", "body", m.ListID); err != nil {
		return err
	}

	if err := validate.FormatOf("list_id", "body", "uuid", m.ListID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *List) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *List) validateUserID(formats strfmt.Registry) error {

	if err := validate.Required("user_id", "body", m.UserID); err != nil {
		return err
	}

	if err := validate.FormatOf("user_id", "body", "uuid", m.UserID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this list based on the context it is used
func (m *List) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBooks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *List) contextValidateBooks(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Books); i++ {

		if m.Books[i] != nil {
			if err := m.Books[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("books" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("books" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *List) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *List) UnmarshalBinary(b []byte) error {
	var res List
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
