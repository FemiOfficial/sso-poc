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
	appIdentityProviderRepository *AppIdentityProviderRepository) (*entitities.App, error) {
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

	app := &entitities.App{
		Name:           request.Name,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		RedirectURI:    request.RedirectURI,
		Live:           false,
		Scopes:         request.Scopes,
		MfaEnabled:     request.MfaEnabled,
		OrganizationID: organizationID,
	}

	err = tx.Create(app).Error
	if err != nil {
		return nil, err
	}

	appIdentityProviders := []*entitities.AppIdentityProvider{}
	if len(request.IdentityProviders) == 0 {
		appIdentityProviders, err = appIdentityProviderRepository.FindAllByFilter(AppIdentityProviderFilter{
			Status:    "active",
			IsDefault: true,
		}, tx)
		if err != nil {
			return nil, err
		}
	} else {
		appIdentityProviders, err = r.getAppProvidersFromRequest(request, app, appIdentityProviderRepository, tx)
		if err != nil {
			return nil, err
		}
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

	app := &entitities.App{}
	return app, query.First(app).Error
}

func (r *AppRepository) getAppProvidersFromRequest(
	request *appTypes.CreateAppRequest,
	app *entitities.App,
	appIdentityProviderRepository *AppIdentityProviderRepository,
	tx *gorm.DB) ([]*entitities.AppIdentityProvider, error) {

	appIdentityProviders := []*entitities.AppIdentityProvider{}
	identityProviderIds := []string{}
	for _, identityProvider := range request.IdentityProviders {
		identityProviderIds = append(identityProviderIds, identityProvider.ID)
	}
	dbIdentityProviders, err := appIdentityProviderRepository.FindAllByFilter(AppIdentityProviderFilter{
		IDs: identityProviderIds,
	}, tx)

	if err != nil {
		return nil, err
	}

	for _, dbIdentityProvider := range dbIdentityProviders {
		requestIdentityProviderIndex := slices.IndexFunc(request.IdentityProviders, func(identityProvider appTypes.IdentityProviders) bool {
			return identityProvider.ID == dbIdentityProvider.IdentityProviderID
		})
		requestIdentityProvider := request.IdentityProviders[requestIdentityProviderIndex]

		scopes := []string{}
		if requestIdentityProvider.Scopes != nil {
			scopes = requestIdentityProvider.Scopes
		} else {
			scopes = dbIdentityProvider.Scopes
		}

		appIdentityProviders = append(appIdentityProviders, &entitities.AppIdentityProvider{
			AppID:              app.ID,
			IdentityProviderID: dbIdentityProvider.IdentityProviderID,
			Status:             "active",
			Scopes:             scopes,
		})
	}
	return appIdentityProviders, nil
}
