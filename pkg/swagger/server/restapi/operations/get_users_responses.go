// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/dccn-tg/dr-gateway/pkg/swagger/server/models"
)

// GetUsersOKCode is the HTTP code returned for type GetUsersOK
const GetUsersOKCode int = 200

/*
GetUsersOK success

swagger:response getUsersOK
*/
type GetUsersOK struct {

	/*
	  In: Body
	*/
	Payload *models.ResponseBodyUsers `json:"body,omitempty"`
}

// NewGetUsersOK creates GetUsersOK with default headers values
func NewGetUsersOK() *GetUsersOK {

	return &GetUsersOK{}
}

// WithPayload adds the payload to the get users o k response
func (o *GetUsersOK) WithPayload(payload *models.ResponseBodyUsers) *GetUsersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users o k response
func (o *GetUsersOK) SetPayload(payload *models.ResponseBodyUsers) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUsersInternalServerErrorCode is the HTTP code returned for type GetUsersInternalServerError
const GetUsersInternalServerErrorCode int = 500

/*
GetUsersInternalServerError failure

swagger:response getUsersInternalServerError
*/
type GetUsersInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ResponseBody500 `json:"body,omitempty"`
}

// NewGetUsersInternalServerError creates GetUsersInternalServerError with default headers values
func NewGetUsersInternalServerError() *GetUsersInternalServerError {

	return &GetUsersInternalServerError{}
}

// WithPayload adds the payload to the get users internal server error response
func (o *GetUsersInternalServerError) WithPayload(payload *models.ResponseBody500) *GetUsersInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users internal server error response
func (o *GetUsersInternalServerError) SetPayload(payload *models.ResponseBody500) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
