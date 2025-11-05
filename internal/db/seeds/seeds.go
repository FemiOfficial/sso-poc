package seeds

import (
	"sso-poc/internal/db"
)

type Seeder struct {
	db *db.Database
}

func NewSeeder(db *db.Database) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Seed() error {

	database := db.InitializeDB()
	defer database.Close()

	seeder := NewSeeder(database)

	err := seeder.SeedIdentityProviders()
	if err != nil {
		return err
	}
	return nil
}
