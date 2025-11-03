package main

import (
	"log"
	"os"
	"sso-poc/internal/db"
	"sso-poc/internal/db/seeds"
)

func main() {

	if os.Getenv("APP_ENV") == "production" {
		log.Fatalf("Seeding is not allowed in production")
		return
	}

	err := seeds.NewSeeder(db.InitializeDB()).Seed()
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
}
