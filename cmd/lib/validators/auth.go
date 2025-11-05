package validators

import (
	"errors"
	"net/http"
	"sso-poc/internal/db/entitities"
	"sso-poc/internal/db/repositories"
)

func ClientAuthValidator(clientId string, clientSecret string, appRepository *repositories.AppRepository) (app *entitities.App, statusCode int, err error) {
	if clientId == "" || clientSecret == "" {
		return nil, http.StatusBadRequest, errors.New("client id and client secret are required")
	}

	app, err = appRepository.FindOneByFilter(repositories.AppFilter{ClientID: clientId}, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if app == nil {
		return nil, http.StatusUnauthorized, errors.New("invalid client id")
	}

	if app.ClientSecret != clientSecret {
		return nil, http.StatusUnauthorized, errors.New("invalid client secret")
	}

	return app, http.StatusOK, nil
}
