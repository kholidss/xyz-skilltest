package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type Config struct {
	AppConfig      `mapstructure:",squash"`
	LoggerConfig   `mapstructure:",squash"`
	DatabaseConfig `mapstructure:",squash"`
	MongoDBConfig  `mapstructure:",squash"`
	BrokerConfig   `mapstructure:",squash"`
	GCPConfig      `mapstructure:",squash"`
	NoSleepConfig  `mapstructure:",squash"`
	StorageConfig  `mapstructure:",squash"`
	CDNConfig      `mapstructure:",squash"`
}

func LoadAllConfigs() (*Config, error) {
	var cnf Config
	err := loadConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cnf)
	if err != nil {
		return nil, err
	}

	return &cnf, nil
}

// FiberConfig func for configuration Fiber app.
func (cnf *Config) FiberConfig() fiber.Config {
	// Return Fiber configuration.
	return fiber.Config{
		AppName:       cnf.AppName,
		StrictRouting: false,
		CaseSensitive: false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			switch code {
			// Not found handle
			case http.StatusNotFound:
				c.Status(http.StatusNotFound).JSON(map[string]interface{}{
					"message": "Sorry, the resource not found!",
				})
			// Method not allowed handle
			case http.StatusMethodNotAllowed:
				c.Status(http.StatusMethodNotAllowed).JSON(map[string]interface{}{
					"message": "Method not allowed!",
				})
			// Default internal server error handle
			case http.StatusInternalServerError:
				c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
					"message": "Something went wrong!",
				})
			}
			return nil
		},
	}
}

func loadConfig() error {
	viper.AutomaticEnv()

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	_ = viper.MergeInConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("Config %v was change", e.Name)
	})

	return nil
}
