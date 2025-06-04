package database

import (
	"backend/pkg/errors"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, ErrMissingDBURL
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

var ErrMissingDBURL = errors.New(errors.ErrBadRequest, "missing DATABASE_URL environment variable")
