package db

import (
	"github.com/ambientis-org/hefesto/internal/db/postgres/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrating User
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	_ = db.AutoMigrate(&models.Doctor{})

	return db, nil
}
