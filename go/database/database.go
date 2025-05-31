package database

import (
    "fmt"
    "log"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "mockapi/config"  // Assuming module name is mockapi
    "mockapi/models"  // Assuming module name is mockapi
)

var DB *gorm.DB

// ConnectDB connects to the database using the provided configuration
// and runs AutoMigrate.
func ConnectDB(cfg config.Config) {
    var err error

    // MySQL DSN format: user:password@tcp(host:port)/dbname?parseTime=true
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    log.Println("✅ Database connection established.")

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
    log.Println("✅ Database migrated successfully.")
}

// GetDB returns the current database instance.
func GetDB() *gorm.DB {
    return DB
}
