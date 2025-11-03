package validators

import (
	"errors"
	"sso-poc/internal/db/repositories"
)

func ValidatorProviders(providers []string, identityProviderRepository *repositories.IdentityProviderRepository) (validProviders []*string, invalidProviders []*string, err error) {
	if len(providers) == 0 {
		return nil, nil, errors.New("providers are required")
	}

	dbProviders, err := identityProviderRepository.FindAllByFilter(repositories.IdentityProviderFilter{
		Names:  providers,
		Status: "active",
	}, nil)
	if err != nil {
		return nil, nil, err
	}

	name := make(map[string]bool, len(dbProviders))
	for _, provider := range dbProviders {
		if provider != nil {
			name[provider.Name] = true
		}
	}

	for _, provider := range providers {
		if name[provider] {
			validProviders = append(validProviders, &provider)
		} else {
			invalidProviders = append(invalidProviders, &provider)
		}
	}

	return validProviders, invalidProviders, nil

}
