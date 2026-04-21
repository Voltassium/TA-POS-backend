package config

import (
	"errors"
	"fmt"
	"log"
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

		viper.AutomaticEnv()

		if err = viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				err = fmt.Errorf("config file not found")
				return
			}
			err = fmt.Errorf("error reading config file: %w", err)
			return
		}

		if err = viper.Unmarshal(&config); err != nil {
			return
		}
	})

	if err != nil {
		log.Fatalf("%s", fmt.Sprintf("error loading config: %s", err.Error()))
	}

	return config
}
