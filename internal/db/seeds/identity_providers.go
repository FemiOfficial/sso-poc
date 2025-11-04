package seeds

import (
	"encoding/json"
	"io"
	"os"
	"sso-poc/internal/db/entitities"

	"gorm.io/gorm"
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
		existsErr := tx.Model(&entitities.IdentityProvider{}).Where("name = ?", identityProvider.Name).First(&entitities.IdentityProvider{}).Error
		if existsErr == gorm.ErrRecordNotFound {

			scopes := identityProvider.Scopes
			if len(scopes) == 0 {
				scopes = nil
			}

			err = tx.Create(&entitities.IdentityProvider{
				Name:        identityProvider.Name,
				DisplayName: identityProvider.DisplayName,
				Scopes:      scopes,
				Status:      identityProvider.Status,
			}).Error
			if err != nil {
				return err
			}
		}  
		
		if existsErr == nil {
			credentials := identityProvider.CredentialFields
			if len(credentials) > 0 {
				existsIdentityProvider := &entitities.IdentityProvider{}
				err = tx.Model(&entitities.IdentityProvider{}).Where("name = ?", identityProvider.Name).First(existsIdentityProvider).Error
				if err != nil {
					return err
				}
				existsIdentityProvider.CredentialFields = credentials
				err = tx.Save(existsIdentityProvider).Error
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit().Error
}
