package database

import (
	"fmt"
	"log"
	// "strconv" // Unused import, can be removed

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mockapi/config" // Assuming module name is mockapi
	"mockapi/models" // Assuming module name is mockapi
)

var DB *gorm.DB

// ConnectDB connects to the database using the provided configuration
// and runs AutoMigrate.
func ConnectDB(cfg config.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSslmode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Or logger.Silent for less verbosity
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Auto-migrate models
	err = DB.AutoMigrate(
		&models.Team{},
		&models.Project{},
		&models.ForwardProxy{},
		&models.Url{},
		&models.MockContent{},
		&models.RequestLog{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database migrated successfully.")
}

// GetDB returns the current database instance.
// This is a simple way to provide access to the DB instance.
// Consider dependency injection for larger applications.
func GetDB() *gorm.DB {
	return DB
}
