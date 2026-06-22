package config

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const configFilePath = "./pkg/config/files"

type Config struct {
	Application    Application    `mapstructure:"application"`
	Authentication Authentication `mapstructure:"authentication"`
	Database       Database       `mapstructure:"database"`
	Logger         Logger         `mapstructure:"logger"`
	ObjectStorage  ObjectStorage  `mapstructure:"object_storage"`
}

var once sync.Once
var config Config

func LoadConfig() Config {
	var (
		err error
	)
	once.Do(func() {
		viper.SetConfigName("env")          // name of config file (without extension)
		viper.SetConfigType("yaml")         // optional, set if the config file is not .json, .yaml, etc.
		viper.AddConfigPath(configFilePath) // optionally look for config in the working directory

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		bindEnvKeys()

		if err = viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				log.Println("Config file not found. Loading from environment variables...")
				err = nil // Ignore error if file is not found
			} else {
				err = fmt.Errorf("error reading config file: %w", err)
				return
			}
		}

		if err = viper.Unmarshal(&config); err != nil {
			err = fmt.Errorf("error unmarshaling config: %w", err)
			return
		}
	})

	if err != nil {
		log.Fatalf("%s", fmt.Sprintf("error loading config: %s", err.Error()))
	}

	return config
}

// bindEnvKeys explicitly binds environment variable names to viper keys.
func bindEnvKeys() {
	bindings := map[string]string{
		// Application
		"APPLICATION_PORT":        "application.port",
		"APPLICATION_ENVIRONMENT": "application.environment",

		// Database
		"DATABASE_HOST":               "database.host",
		"DATABASE_PORT":               "database.port",
		"DATABASE_USER":               "database.user",
		"DATABASE_PASSWORD":           "database.password",
		"DATABASE_NAME":               "database.name",
		"DATABASE_SSL_MODE":           "database.ssl_mode",
		"DATABASE_MAX_OPEN_CONN":      "database.max_open_conn",
		"DATABASE_MAX_OPEN_IDLE_CONN": "database.max_open_idle_conn",
		"DATABASE_MAX_IDLE_CONN":      "database.max_idle_conn",
		"DATABASE_URL":                "database.url",

		// Authentication
		"AUTHENTICATION_ENCRYPT_KEY":          "authentication.encrypt_key",
		"AUTHENTICATION_ACCESS_SECRET_KEY":    "authentication.access_secret_key",
		"AUTHENTICATION_REFRESH_SECRET_KEY":   "authentication.refresh_secret_key",
		"AUTHENTICATION_ISSUER":               "authentication.issuer",
		"AUTHENTICATION_ACCESS_TOKEN_EXPIRY":  "authentication.access_token_expiry",
		"AUTHENTICATION_REFRESH_TOKEN_EXPIRY": "authentication.refresh_token_expiry",

		// Logger
		"LOGGER_ENVIRONMENT": "logger.environment",
		"LOGGER_LOG_LEVEL":   "logger.log_level",
		"LOGGER_ENCODING":    "logger.encoding",

		// Object Storage (MinIO)
		"OBJECT_STORAGE_ENDPOINT":           "object_storage.endpoint",
		"OBJECT_STORAGE_BUCKET":             "object_storage.bucket",
		"OBJECT_STORAGE_ACCESS_KEY":         "object_storage.access_key",
		"OBJECT_STORAGE_SECRET_KEY":         "object_storage.secret_key",
		"OBJECT_STORAGE_USE_SSL":            "object_storage.use_ssl",
		"OBJECT_STORAGE_PRESIGN_EXPIRATION": "object_storage.presign_expiration",
		"OBJECT_STORAGE_MAX_FILE_SIZE":      "object_storage.max_file_size",
	}

	for envKey, viperKey := range bindings {
		_ = viper.BindEnv(viperKey, envKey)
	}
}
