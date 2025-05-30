package config

import (
	"fmt"
	"go-gin-gorm-api/internal/models" // Adjusted import path
	"strings"                         // Added

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres" // Added
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and performs auto-migration.
// It now uses GetDBDSN() from the AppConfig.
func InitDB() (*gorm.DB, error) {
	if AppConfig == nil {
		return nil, fmt.Errorf("configuration not loaded. Call LoadConfig first")
	}

	dsn := GetDBDSN()
	if dsn == "" {
		return nil, fmt.Errorf("generated DSN is empty. Check DB_DRIVER in config")
	}

	var dialector gorm.Dialector
	switch strings.ToLower(AppConfig.DBDriver) {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres", "postgresql":
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER: %s", AppConfig.DBDriver)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database (driver: %s): %w", AppConfig.DBDriver, err)
	}

	// Auto-migrate models
	fmt.Println("Running auto-migrations...")
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
	fmt.Println("Auto-migrations completed.")
	return DB, nil
}
