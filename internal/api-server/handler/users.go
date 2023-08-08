package handler

import (
	"strings"

	"github.com/dccn-tg/dr-gateway/pkg/dr"
	"github.com/dccn-tg/dr-gateway/pkg/swagger/server/models"
	"github.com/dccn-tg/dr-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// GetUsers returns all users.
func GetUsers(ucache *UsersCache) func(params operations.GetUsersParams) middleware.Responder {
	return func(params operations.GetUsersParams) middleware.Responder {

		users := []*models.ResponseBodyUserMetadata{}

		for _, u := range ucache.GetUsers() {
			users = append(users, makeResponseBodyUserMetadata(u))
		}

		return operations.NewGetUsersOK().WithPayload(&models.ResponseBodyUsers{
			Users: users,
		})
	}
}

// GetUsersOfOu returns all users of an organisational unit.
func GetUsersOfOu(ucache *UsersCache) func(params operations.GetUsersOuIDParams) middleware.Responder {
	return func(params operations.GetUsersOuIDParams) middleware.Responder {
		id := strings.ToUpper(params.ID)

		users := []*models.ResponseBodyUserMetadata{}
		for _, u := range ucache.GetUsers() {
			if contains(u.OrganizationalUnits, id) {
				users = append(users, makeResponseBodyUserMetadata(u))
			}
		}

		return operations.NewGetUsersOuIDOK().WithPayload(&models.ResponseBodyUsers{
			Users: users,
		})
	}
}

// SearchUsers returns all users with matched email and/or display name.
func SearchUsers(ucache *UsersCache) func(params operations.GetUsersSearchParams) middleware.Responder {
	return func(params operations.GetUsersSearchParams) middleware.Responder {

		users := []*models.ResponseBodyUserMetadata{}

		// return immediately with an empty slice if both email and name query values are not specified
		if params.Email == nil && params.Name == nil {
			return operations.NewGetUsersOuIDOK().WithPayload(&models.ResponseBodyUsers{
				Users: users,
			})
		}

		for _, u := range ucache.GetUsers() {

			// check email
			if params.Email != nil && u.Email != *params.Email {
				continue
			}

			// check display name
			if params.Name != nil && !strings.Contains(u.DisplayName, *params.Name) {
				continue
			}

			// this user has matched email and/or display name
			users = append(users, makeResponseBodyUserMetadata(u))
		}

		return operations.NewGetUsersOuIDOK().WithPayload(&models.ResponseBodyUsers{
			Users: users,
		})
	}
}

func makeResponseBodyUserMetadata(u *dr.DRUser) *models.ResponseBodyUserMetadata {
	return &models.ResponseBodyUserMetadata{
		DisplayName:         &u.DisplayName,
		Email:               strfmt.Email(u.Email),
		IdentityProvider:    &u.IdentityProvider,
		OrganisationalUnits: u.OrganizationalUnits,
	}
}
