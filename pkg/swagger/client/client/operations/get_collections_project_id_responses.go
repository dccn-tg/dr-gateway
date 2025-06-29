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

// GetCollectionsProjectIDReader is a Reader for the GetCollectionsProjectID structure.
type GetCollectionsProjectIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCollectionsProjectIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCollectionsProjectIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetCollectionsProjectIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /collections/project/{id}] GetCollectionsProjectID", response, response.Code())
	}
}

// NewGetCollectionsProjectIDOK creates a GetCollectionsProjectIDOK with default headers values
func NewGetCollectionsProjectIDOK() *GetCollectionsProjectIDOK {
	return &GetCollectionsProjectIDOK{}
}

/*
GetCollectionsProjectIDOK describes a response with status code 200, with default header values.

success
*/
type GetCollectionsProjectIDOK struct {
	Payload *models.ResponseBodyCollections
}

// IsSuccess returns true when this get collections project Id o k response has a 2xx status code
func (o *GetCollectionsProjectIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get collections project Id o k response has a 3xx status code
func (o *GetCollectionsProjectIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get collections project Id o k response has a 4xx status code
func (o *GetCollectionsProjectIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get collections project Id o k response has a 5xx status code
func (o *GetCollectionsProjectIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get collections project Id o k response a status code equal to that given
func (o *GetCollectionsProjectIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get collections project Id o k response
func (o *GetCollectionsProjectIDOK) Code() int {
	return 200
}

func (o *GetCollectionsProjectIDOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /collections/project/{id}][%d] getCollectionsProjectIdOK %s", 200, payload)
}

func (o *GetCollectionsProjectIDOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /collections/project/{id}][%d] getCollectionsProjectIdOK %s", 200, payload)
}

func (o *GetCollectionsProjectIDOK) GetPayload() *models.ResponseBodyCollections {
	return o.Payload
}

func (o *GetCollectionsProjectIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBodyCollections)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCollectionsProjectIDInternalServerError creates a GetCollectionsProjectIDInternalServerError with default headers values
func NewGetCollectionsProjectIDInternalServerError() *GetCollectionsProjectIDInternalServerError {
	return &GetCollectionsProjectIDInternalServerError{}
}

/*
GetCollectionsProjectIDInternalServerError describes a response with status code 500, with default header values.

failure
*/
type GetCollectionsProjectIDInternalServerError struct {
	Payload *models.ResponseBody500
}

// IsSuccess returns true when this get collections project Id internal server error response has a 2xx status code
func (o *GetCollectionsProjectIDInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get collections project Id internal server error response has a 3xx status code
func (o *GetCollectionsProjectIDInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get collections project Id internal server error response has a 4xx status code
func (o *GetCollectionsProjectIDInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get collections project Id internal server error response has a 5xx status code
func (o *GetCollectionsProjectIDInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get collections project Id internal server error response a status code equal to that given
func (o *GetCollectionsProjectIDInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get collections project Id internal server error response
func (o *GetCollectionsProjectIDInternalServerError) Code() int {
	return 500
}

func (o *GetCollectionsProjectIDInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /collections/project/{id}][%d] getCollectionsProjectIdInternalServerError %s", 500, payload)
}

func (o *GetCollectionsProjectIDInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /collections/project/{id}][%d] getCollectionsProjectIdInternalServerError %s", 500, payload)
}

func (o *GetCollectionsProjectIDInternalServerError) GetPayload() *models.ResponseBody500 {
	return o.Payload
}

func (o *GetCollectionsProjectIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ResponseBody500)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
