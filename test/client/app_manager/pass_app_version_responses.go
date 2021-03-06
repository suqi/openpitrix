// Code generated by go-swagger; DO NOT EDIT.

package app_manager

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"openpitrix.io/openpitrix/test/models"
)

// PassAppVersionReader is a Reader for the PassAppVersion structure.
type PassAppVersionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PassAppVersionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPassAppVersionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPassAppVersionOK creates a PassAppVersionOK with default headers values
func NewPassAppVersionOK() *PassAppVersionOK {
	return &PassAppVersionOK{}
}

/*PassAppVersionOK handles this case with default header values.

A successful response.
*/
type PassAppVersionOK struct {
	Payload *models.OpenpitrixPassAppVersionResponse
}

func (o *PassAppVersionOK) Error() string {
	return fmt.Sprintf("[POST /v1/app_version/action/{role}/pass][%d] passAppVersionOK  %+v", 200, o.Payload)
}

func (o *PassAppVersionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OpenpitrixPassAppVersionResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
