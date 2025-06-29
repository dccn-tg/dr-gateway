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

// GetUsersReader is a Reader for the GetUsers structure.
type GetUsersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetUsersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetUsersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetUsersInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /users] GetUsers", response, response.Code())
	}
}

// NewGetUsersOK creates a GetUsersOK with default headers values
func NewGetUsersOK() *GetUsersOK {
	return &GetUsersOK{}
}

/*
GetUsersOK describes a response with status code 200, with default header values.

success
*/
type GetUsersOK struct {
	Payload *models.ResponseBodyUsers
}

// IsSuccess returns true when this get users o k response has a 2xx status code
func (o *GetUsersOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get users o k response has a 3xx status code
func (o *GetUsersOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get users o k response has a 4xx status code
func (o *GetUsersOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get users o k response has a 5xx status code
func (o *GetUsersOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get users o k response a status code equal to that given
func (o *GetUsersOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get users o k response
func (o *GetUsersOK) Code() int {
	return 200
}

func (o *GetUsersOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /users][%d] getUsersOK %s", 200, payload)
}

func (o *GetUsersOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /users][%d] getUsersOK %s", 200, payload)
}

func (o *GetUsersOK) GetPayload() *models.ResponseBodyUsers {
	return o.Payload
}

func (o *GetUsersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyUsers)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetUsersInternalServerError creates a GetUsersInternalServerError with default headers values
func NewGetUsersInternalServerError() *GetUsersInternalServerError {
	return &GetUsersInternalServerError{}
}

/*
GetUsersInternalServerError describes a response with status code 500, with default header values.

failure
*/
type GetUsersInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this get users internal server error response has a 2xx status code
func (o *GetUsersInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get users internal server error response has a 3xx status code
func (o *GetUsersInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get users internal server error response has a 4xx status code
func (o *GetUsersInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get users internal server error response has a 5xx status code
func (o *GetUsersInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get users internal server error response a status code equal to that given
func (o *GetUsersInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get users internal server error response
func (o *GetUsersInternalServerError) Code() int {
	return 500
}

func (o *GetUsersInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /users][%d] getUsersInternalServerError %s", 500, payload)
}

func (o *GetUsersInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /users][%d] getUsersInternalServerError %s", 500, payload)
}

func (o *GetUsersInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *GetUsersInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
