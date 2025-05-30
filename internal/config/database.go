package config

import (
	"fmt"
	"go-gin-gorm-api/internal/models" // Adjusted import path

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and performs auto-migration.
func InitDB(dsn string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(
		&models.Team{},
		&models.Project{},
		&models.Url{},
		&models.MockContent{},
		&models.RequestLog{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	return DB, nil
}
