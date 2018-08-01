// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// OpenpitrixCreateAppVersionRequest openpitrix create app version request
// swagger:model openpitrixCreateAppVersionRequest
type OpenpitrixCreateAppVersionRequest struct {

	// app id
	AppID string `json:"app_id,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// owner
	Owner string `json:"owner,omitempty"`

	// package name
	PackageName string `json:"package_name,omitempty"`

	// sequence
	Sequence int64 `json:"sequence,omitempty"`
}

// Validate validates this openpitrix create app version request
func (m *OpenpitrixCreateAppVersionRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *OpenpitrixCreateAppVersionRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OpenpitrixCreateAppVersionRequest) UnmarshalBinary(b []byte) error {
	var res OpenpitrixCreateAppVersionRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
