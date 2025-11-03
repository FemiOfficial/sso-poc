package validators

import (
	"errors"
	"net/http"
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
)

func ClientAuthValidator(clientId string, clientSecret string, db *gorm.DB) (app *entitities.App, statusCode int, err error) {
	if clientId == "" || clientSecret == "" {
		return nil, http.StatusBadRequest, errors.New("client id and client secret are required")
	}

	app = &entitities.App{}
	err = db.Where("client_id = ?", clientId).First(app).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if app.ClientSecret != clientSecret {
		return nil, http.StatusUnauthorized, errors.New("invalid client secret")
	}

	return app, http.StatusOK, nil
}
