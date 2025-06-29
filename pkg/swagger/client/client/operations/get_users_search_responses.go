// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/dccn-tg/dr-gateway/pkg/swagger/client/models"
)

// GetUsersSearchReader is a Reader for the GetUsersSearch structure.
type GetUsersSearchReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetUsersSearchReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetUsersSearchOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetUsersSearchInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /users/search] GetUsersSearch", response, response.Code())
	}
}

// NewGetUsersSearchOK creates a GetUsersSearchOK with default headers values
func NewGetUsersSearchOK() *GetUsersSearchOK {
	return &GetUsersSearchOK{}
}

/*
GetUsersSearchOK describes a response with status code 200, with default header values.

success
*/
type GetUsersSearchOK struct {
	Payload *models.ResponseBodyUsers
}

// IsSuccess returns true when this get users search o k response has a 2xx status code
func (o *GetUsersSearchOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get users search o k response has a 3xx status code
func (o *GetUsersSearchOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get users search o k response has a 4xx status code
func (o *GetUsersSearchOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get users search o k response has a 5xx status code
func (o *GetUsersSearchOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get users search o k response a status code equal to that given
func (o *GetUsersSearchOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get users search o k response
func (o *GetUsersSearchOK) Code() int {
	return 200
}

func (o *GetUsersSearchOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /users/search][%d] getUsersSearchOK %s", 200, payload)
}

func (o *GetUsersSearchOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /users/search][%d] getUsersSearchOK %s", 200, payload)
}

func (o *GetUsersSearchOK) GetPayload() *models.ResponseBodyUsers {
	return o.Payload
}

func (o *GetUsersSearchOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyUsers)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetUsersSearchInternalServerError creates a GetUsersSearchInternalServerError with default headers values
func NewGetUsersSearchInternalServerError() *GetUsersSearchInternalServerError {
	return &GetUsersSearchInternalServerError{}
}

/*
GetUsersSearchInternalServerError describes a response with status code 500, with default header values.

failure
*/
type GetUsersSearchInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this get users search internal server error response has a 2xx status code
func (o *GetUsersSearchInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get users search internal server error response has a 3xx status code
func (o *GetUsersSearchInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get users search internal server error response has a 4xx status code
func (o *GetUsersSearchInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get users search internal server error response has a 5xx status code
func (o *GetUsersSearchInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get users search internal server error response a status code equal to that given
func (o *GetUsersSearchInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get users search internal server error response
func (o *GetUsersSearchInternalServerError) Code() int {
	return 500
}

func (o *GetUsersSearchInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /users/search][%d] getUsersSearchInternalServerError %s", 500, payload)
}

func (o *GetUsersSearchInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /users/search][%d] getUsersSearchInternalServerError %s", 500, payload)
}

func (o *GetUsersSearchInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *GetUsersSearchInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
