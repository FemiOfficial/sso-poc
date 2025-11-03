package misc

import (
	miscTypes "sso-poc/cmd/api/server/dashboard/misc/types"
	"sso-poc/internal/db"
	"sso-poc/internal/db/entitities"
	"sso-poc/internal/db/repositories"
	"strings"

	"github.com/gin-gonic/gin"
)

type MiscService struct {
	db                         *db.Database
	identityProviderRepository *repositories.IdentityProviderRepository
}

func CreateMiscService(db *db.Database) *MiscService {
	return &MiscService{db: db, identityProviderRepository: repositories.CreateIdentityProviderRepository(db.DB)}
}

func (s *MiscService) GetIdentityProviders(ctx *gin.Context) (*miscTypes.GetIDPResponse, error) {
	var getIDPRequest miscTypes.GetIDPRequest = ctx.MustGet("request").(miscTypes.GetIDPRequest)

	ids := []string{}
	if getIDPRequest.IDs != "" {
		ids = strings.Split(getIDPRequest.IDs, ",")
	}

	list, err := s.identityProviderRepository.FindAllByFilter(repositories.IdentityProviderFilter{
		Status: getIDPRequest.Status,
		IDs:    ids,
		Name:   getIDPRequest.Name,
		Scopes: getIDPRequest.Scopes,
	}, nil)
	if err != nil {
		return nil, err
	}

	out := make([]entitities.IdentityProvider, 0, len(list))
	for _, p := range list {
		if p != nil {
			out = append(out, *p)
		}
	}

	return &miscTypes.GetIDPResponse{IdentityProviders: out}, nil
}
