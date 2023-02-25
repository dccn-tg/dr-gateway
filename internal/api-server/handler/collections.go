package handler

import (
	"github.com/Donders-Institute/dr-gateway/pkg/dr"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

// GetCollections returns all collections.
func GetCollections(ccache *CollectionsCache) func(params operations.GetCollectionsParams) middleware.Responder {
	return func(params operations.GetCollectionsParams) middleware.Responder {

		colls := []*models.ResponseBodyCollectionMetadata{}

		for _, c := range ccache.GetCollections() {
			colls = append(colls, &models.ResponseBodyCollectionMetadata{
				Identifier:         &c.Identifier,
				IdentifierDOI:      &c.IdentifierDOI,
				Path:               &c.Path,
				OrganizationalUnit: &c.OrganisationalUnit,
				QuotaInBytes:       &c.QuotaInBytes,
				NumberOfFiles:      &c.NumberOfFiles,
				SizeInBytes:        &c.SizeInBytes,
				State:              collectionState(c.State),
				Type:               collectionType(c.Type),
			})
		}

		return operations.NewGetCollectionsOK().WithPayload(&models.ResponseBodyCollections{
			Collections: colls,
		})
	}
}

func collectionState(s dr.CollectionState) *models.CollectionState {
	switch s {
	case dr.Editable:
		return models.CollectionStateEditable.Pointer()
	case dr.ReviewableInternal:
		return models.CollectionStateReviewableInternal.Pointer()
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
