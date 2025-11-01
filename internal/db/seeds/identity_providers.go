package seeds

import (
	"encoding/json"
	"io"
	"os"
	"sso-poc/internal/db/entitities"
)

func (s *Seeder) SeedIdentityProviders() error {

	jsonFile, err := os.Open("internal/db/seeds/data/idp.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var identityProviders []entitities.IdentityProvider
	err = json.Unmarshal(byteValue, &identityProviders)
	if err != nil {
		return err
	}

	tx := s.db.DB.Begin()
	defer tx.Rollback()

	for _, identityProvider := range identityProviders {
		err = tx.Create(&entitities.IdentityProvider{
			Name:        identityProvider.Name,
			DisplayName: identityProvider.DisplayName,
			Scopes:      identityProvider.Scopes,
			Status:      identityProvider.Status,
		}).Error
		if err != nil {
			return err
		}
	}

	return tx.Commit().Error
}
