package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
// Values are read from config file and environment variables.
type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSslMode  string `mapstructure:"DB_SSL_MODE"` // e.g. "disable" for local, "require" for prod
	AppBaseURL string `mapstructure:"APP_BASE_URL"`

	// Pusher credentials (can be added later if Pusher is fully integrated)
	// PusherAppID  string `mapstructure:"PUSHER_APP_ID"`
	// PusherKey    string `mapstructure:"PUSHER_KEY"`
	// PusherSecret string `mapstructure:"PUSHER_SECRET"`
	// PusherCluster string `mapstructure:"PUSHER_CLUSTER"`
}

// AppConfig is the global configuration object.
var AppConfig *Config

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path) // e.g., "." for root directory
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") // or json, toml

	viper.AutomaticEnv()                                         // Read matching environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // For nested env vars e.g. SERVER.PORT -> SERVER_PORT

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if this is intentional (rely on env vars)
			fmt.Println("Config file not found, relying on environment variables.")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	AppConfig = &Config{} // Initialize AppConfig before unmarshalling
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// You might want to add validation for critical config values here
	if AppConfig.ServerPort == "" {
		AppConfig.ServerPort = "8080" // Default port if not set
	}
	if AppConfig.DBDriver == "" {
		// Or return an error if a critical value is missing
		fmt.Println("DB_DRIVER not set, defaulting to mysql for DSN generation if applicable.")
		// AppConfig.DBDriver = "mysql" // Example default
	}


	return AppConfig, nil
}

// GetDBDSN constructs the database DSN string from AppConfig.
// It supports mysql and postgres drivers.
func GetDBDSN() string {
	if AppConfig == nil {
		// This should not happen if LoadConfig is called first.
		// Consider panicking or returning an error.
		fmt.Println("Error: AppConfig is not loaded. Call LoadConfig first.")
		return "" // Or panic
	}

	switch strings.ToLower(AppConfig.DBDriver) {
	case "mysql":
		// Example: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			AppConfig.DBUser,
			AppConfig.DBPassword,
			AppConfig.DBHost,
			AppConfig.DBPort,
			AppConfig.DBName,
		)
	case "postgres", "postgresql":
		// Example: "host=localhost port=5432 user=youruser password=yourpassword dbname=yourdbname sslmode=disable"
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			AppConfig.DBHost,
			AppConfig.DBPort,
			AppConfig.DBUser,
			AppConfig.DBPassword,
			AppConfig.DBName,
			AppConfig.DBSslMode,
		)
	default:
		fmt.Printf("Unsupported DB_DRIVER: %s. Returning empty DSN.\n", AppConfig.DBDriver)
		return ""
	}
}
