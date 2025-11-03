package repositories

import (
	appTypes "sso-poc/cmd/api/server/dashboard/app/types"
	entitities "sso-poc/internal/db/entitities"
	"sso-poc/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type AppRepository struct {
	db *gorm.DB
}

type AppFilter struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	ClientID       string `json:"client_id"`
}

func CreateAppRepository(db *gorm.DB) *AppRepository {
	return &AppRepository{db: db}
}

func (r *AppRepository) Create(
	request *appTypes.CreateAppRequest,
	organizationID string, tx *gorm.DB,
	appIdentityProviderRepository *AppIdentityProviderRepository,
	identityProviderRepository *IdentityProviderRepository,
) (*entitities.App, error) {
	if tx == nil {
		tx = r.db
	}

	clientID := uuid.New().String()
	clientSecret, err := utils.GenerateRandomString(128)
	if err != nil {
		return nil, err
	}

	if len(request.Scopes) == 0 {
		request.Scopes = []entitities.AppScope{entitities.AuthScope, entitities.AuditLogScope, entitities.UserManagementScope}
	}

	scopes := []string{}
	for _, scope := range request.Scopes {
		scopes = append(scopes, string(scope))
	}

	app := &entitities.App{
		Name:           request.Name,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		RedirectURI:    request.RedirectURI,
		Live:           false,
		Scopes:         scopes,
		MfaEnabled:     request.MfaEnabled,
		OrganizationID: organizationID,
	}

	err = tx.Create(app).Error
	if err != nil {
		return nil, err
	}

	identityProviders := []*entitities.IdentityProvider{}
	if len(request.IdentityProviders) == 0 {
		identityProviders, err = identityProviderRepository.FindAllByFilter(IdentityProviderFilter{
			Status:    "active",
			IsDefault: true,
		}, tx)
		if err != nil {
			return nil, err
		}
	} else {
		identityProviders, err = r.getProvidersFromRequest(request, identityProviderRepository, tx)
		if err != nil {
			return nil, err
		}
	}

	appIdentityProviders := []*entitities.AppIdentityProvider{}
	for _, identityProvider := range identityProviders {
		appIdentityProviders = append(appIdentityProviders, &entitities.AppIdentityProvider{
			AppID:      app.ID,
			IdentityProviderID: identityProvider.ID,
			Status:             "active",
			Scopes:             identityProvider.Scopes,
		})
	}


	err = appIdentityProviderRepository.CreateMany(appIdentityProviders, tx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *AppRepository) FindOneByFilter(filter AppFilter, tx *gorm.DB) (*entitities.App, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.Model(&entitities.App{})

	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.OrganizationID != "" {
		query = query.Where("organization_id = ?", filter.OrganizationID)
	}

	if filter.ClientID != "" {
		query = query.Where("client_id = ?", filter.ClientID)
	}

	app := &entitities.App{}
	return app, query.First(app).Error
}

func (r *AppRepository) getProvidersFromRequest(
	request *appTypes.CreateAppRequest,
	identityProviderRepository *IdentityProviderRepository,
	tx *gorm.DB) ([]*entitities.IdentityProvider, error) {

	appIdentityProviders := []*entitities.IdentityProvider{}
	identityProviderIds := []string{}
	for _, identityProvider := range request.IdentityProviders {
		identityProviderIds = append(identityProviderIds, identityProvider.ID)
	}
	dbIdentityProviders, err := identityProviderRepository.FindAllByFilter(IdentityProviderFilter{
		IDs: identityProviderIds,
	}, tx)

	if err != nil {
		return nil, err
	}

	for _, dbIdentityProvider := range dbIdentityProviders {
		requestIdentityProviderIndex := slices.IndexFunc(request.IdentityProviders, func(identityProvider appTypes.IdentityProviders) bool {
			return identityProvider.ID == dbIdentityProvider.ID
		})
		requestIdentityProvider := request.IdentityProviders[requestIdentityProviderIndex]

		scopes := []string{}
		if requestIdentityProvider.Scopes != nil {
			scopes = requestIdentityProvider.Scopes
		} else {
			scopes = dbIdentityProvider.Scopes
		}

		appIdentityProviders = append(appIdentityProviders, &entitities.IdentityProvider{
			BaseEntity: entitities.BaseEntity{ID: dbIdentityProvider.ID},
			Name:       dbIdentityProvider.Name,
			Scopes:     scopes,
			Status:     "active",
			IsDefault:  dbIdentityProvider.IsDefault,
		})
	}
	return appIdentityProviders, nil
}
