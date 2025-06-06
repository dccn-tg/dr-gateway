package handler

import (
	"strings"

	"github.com/dccn-tg/dr-gateway/pkg/dr"
	"github.com/dccn-tg/dr-gateway/pkg/swagger/server/models"
	"github.com/dccn-tg/dr-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

// GetCollections returns all collections.
func GetCollections(ccache *CollectionsCache) func(params operations.GetCollectionsParams) middleware.Responder {
	return func(params operations.GetCollectionsParams) middleware.Responder {

		colls := []*models.ResponseBodyCollectionMetadata{}

		for _, c := range ccache.GetCollections() {
			colls = append(colls, makeResponseBodyCollectionMetadata(c))
		}

		return operations.NewGetCollectionsOK().WithPayload(&models.ResponseBodyCollections{
			Collections: colls,
		})
	}
}

// GetCollectionsOfOu returns all collections of an organisational unit.
func GetCollectionsOfOu(ccache *CollectionsCache) func(params operations.GetCollectionsOuIDParams) middleware.Responder {
	return func(params operations.GetCollectionsOuIDParams) middleware.Responder {
		id := strings.ToLower(params.ID)

		colls := []*models.ResponseBodyCollectionMetadata{}
		for _, c := range ccache.GetCollections() {
			if strings.ToLower(c.OrganisationalUnit) == id {
				colls = append(colls, makeResponseBodyCollectionMetadata(c))
			}
		}

		return operations.NewGetCollectionsOuIDOK().WithPayload(&models.ResponseBodyCollections{
			Collections: colls,
		})
	}
}

// GetCollectionsOfProject returns all collections of a project.
func GetCollectionsOfProject(ccache *CollectionsCache) func(params operations.GetCollectionsProjectIDParams) middleware.Responder {
	return func(params operations.GetCollectionsProjectIDParams) middleware.Responder {
		id := strings.ToLower(params.ID)

		colls := []*models.ResponseBodyCollectionMetadata{}
		for _, c := range ccache.GetCollections() {
			if strings.ToLower(c.ProjectID) == id {
				colls = append(colls, makeResponseBodyCollectionMetadata(c))
			}
		}

		return operations.NewGetCollectionsProjectIDOK().WithPayload(&models.ResponseBodyCollections{
			Collections: colls,
		})
	}
}

func makeResponseBodyCollectionMetadata(c *dr.DRCollection) *models.ResponseBodyCollectionMetadata {
	return &models.ResponseBodyCollectionMetadata{
		Identifier:         &c.Identifier,
		IdentifierDOI:      &c.IdentifierDOI,
		ProjectID:          &c.ProjectID,
		Path:               &c.Path,
		OrganisationalUnit: &c.OrganisationalUnit,
		QuotaInBytes:       &c.QuotaInBytes,
		NumberOfFiles:      &c.NumberOfFiles,
		SizeInBytes:        &c.SizeInBytes,
		State:              collectionState(c.State),
		Type:               collectionType(c.Type),
	}
}

func collectionState(s dr.CollectionState) *models.CollectionState {
	switch s {
	case dr.Editable:
		return models.CollectionStateEditable.Pointer()
	case dr.ReviewableInternal:
		return models.CollectionStateReviewableInternal.Pointer()
	case dr.FairReview:
		return models.CollectionStateFairReview.Pointer()
	case dr.ReviewableExternal:
		return models.CollectionStateReviewableExternal.Pointer()
	case dr.Archived:
		return models.CollectionStateArchived.Pointer()
	case dr.Published:
		return models.CollectionStatePublished.Pointer()
	default:
		return models.CollectionStateUnknown.Pointer()
	}
}

func collectionType(t dr.CollectionType) *models.CollectionType {
	switch t {
	case dr.DAC:
		return models.CollectionTypeDac.Pointer()
	case dr.RDC:
		return models.CollectionTypeRdc.Pointer()
	case dr.DSC:
		return models.CollectionTypeDsc.Pointer()
	default:
		return models.CollectionTypeUnknown.Pointer()
	}
}
