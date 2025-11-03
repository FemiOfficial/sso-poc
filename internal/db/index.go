package db

import (
	"context"
	"log"
	"os"
	_ "github.com/joho/godotenv/autoload"
	"sso-poc/internal/db/entitities"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func InitializeDB() *Database {
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		panic("DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dburl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		panic("failed to ping database")
	}

	log.Println("✅ Database connected successfully")

	err = AutoMigrate(&Database{DB: db})
	if err != nil {
		panic("failed to migrate database")
	}
	return &Database{DB: db}
}

// internal/db/index.go
func AutoMigrate(db *Database) error {
	if os.Getenv("APP_ENV") == "production" {
		log.Println("Skipping AutoMigrate in production")
		return nil
	}

	return db.DB.AutoMigrate(
		&entitities.App{},
		&entitities.Organization{},
		&entitities.IdentityProvider{},
		&entitities.AppIdentityProvider{},
		&entitities.AuthRequest{},
		&entitities.User{},
		&entitities.Vault{},
	)
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) HealthCheck(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
