package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	DBSslmode      string `mapstructure:"DB_SSLMODE"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	JWTSecretKey      string        `mapstructure:"JWT_SECRET_KEY"`
	JWTExpiration     time.Duration `mapstructure:"JWT_EXPIRATION_HOURS"`

	BaseURL       string `mapstructure:"BASE_URL"`
	ServerPort    string `mapstructure:"SERVER_PORT"`

	GlobalMaxAllowedRequests int `mapstructure:"GLOBAL_MAX_ALLOWED_REQUESTS"`
	GlobalTimeWindowSeconds  int `mapstructure:"GLOBAL_TIME_WINDOW_SECONDS"`

	// DefaultTeamID uint   `mapstructure:"DEFAULT_TEAM_ID"` // Uncomment if needed
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path) // Path to look for the config file in
	viper.SetConfigName(".env") // Name of config file (without extension)
	viper.SetConfigType("env")  // Config file type

	viper.AutomaticEnv() // Read in environment variables that match

	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; ignore error if desired
		log.Println("Config file not found, relying on environment variables.")
		// Create a dummy .env file if it doesn't exist to prevent viper error on subsequent calls in some setups
		// Though, for .env type, this might not be strictly necessary if AutomaticEnv is working.
		if _, err := os.Stat(path + "/.env"); os.IsNotExist(err) {
			// Potentially create an empty .env if needed, or ensure viper doesn't error out.
			// For now, we assume AutomaticEnv() is sufficient if file not found.
		}
	} else {
		// Config file was found but another error was produced
		log.Printf("Error reading config file: %s", err)
		return
	}

	// Unmarshal the config into the Config struct
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("Unable to decode into struct, %v", err)
		return
	}

	// Handle JWTExpiration separately as it needs parsing from hours
	jwtExpHoursStr := viper.GetString("JWT_EXPIRATION_HOURS")
	if jwtExpHours, convErr := strconv.Atoi(jwtExpHoursStr); convErr == nil {
		config.JWTExpiration = time.Duration(jwtExpHours) * time.Hour
	} else {
		// Default if not set or invalid
		config.JWTExpiration = 72 * time.Hour
		log.Printf("Invalid or missing JWT_EXPIRATION_HOURS, defaulting to %v", config.JWTExpiration)
	}

    if config.ServerPort == "" {
        config.ServerPort = "8080" // Default port
    }

	return
}
