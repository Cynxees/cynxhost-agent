package dependencies

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name          string `mapstructure:"name"`
		Address       string `mapstructure:"address"`
		PrivateIp     string `mapstructure:"privateIp"`
		Key           string `mapstructure:"key"`
		Port          int    `mapstructure:"port"`
		WebsocketPort int    `mapstructure:"websocketPort"`
		Debug         bool   `mapstructure:"debug"`
	} `mapstructure:"app"`

	Central struct {
		PrivateIp string `mapstructure:"privateIp"`
		PublicIp  string `mapstructure:"publicIp"`
		Port      int    `mapstructure:"port"`
	} `mapstructure:"central"`

	Log struct {
		MinecraftLogFilePath string `mapstructure:"minecraftLogFilePath"`
	}

	Router struct {
		Default string `mapstructure:"default"`
	} `mapstructure:"router"`

	Database struct {
		MySQL struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Database string `mapstructure:"database"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
			Dialect  string `mapstructure:"dialect"`
			Pool     struct {
				Max     int `mapstructure:"max"`
				Min     int `mapstructure:"min"`
				Acquire int `mapstructure:"acquire"`
				Idle    int `mapstructure:"idle"`
			} `mapstructure:"pool"`
		} `mapstructure:"mysql"`

		Redis struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Password string `mapstructure:"password"`
		} `mapstructure:"redis"`

		Elasticsearch struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
			Log  string `mapstructure:"log"`
		} `mapstructure:"elasticsearch"`

		RabbitMQ struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		} `mapstructure:"rabbitmq"`
	} `mapstructure:"database"`

	Logging struct {
		Level        string   `mapstructure:"level"`
		Format       string   `mapstructure:"format"`
		Destinations []string `mapstructure:"destinations"`
	} `mapstructure:"logging"`

	Aws struct {
		AccessKeyId     string `mapstructure:"accessKeyId"`
		AccessKeySecret string `mapstructure:"accessKeySecret"`
	} `mapstructure:"aws"`

	Security struct {
		JWT struct {
			Secret    string `mapstructure:"secret"`
			ExpiresIn string `mapstructure:"expiresIn"`
		} `mapstructure:"jwt"`

		CORS struct {
			Enabled bool   `mapstructure:"enabled"`
			Origin  string `mapstructure:"origin"`
		} `mapstructure:"cors"`
	} `mapstructure:"security"`
}

func LoadConfig(path string) (*Config, error) {
	// Specify the config file path
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Convert `.` to `_` in env vars

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Manually replace placeholders like ${ENV_VAR} with actual environment variable values
	configMap := viper.AllSettings() // Get all config as a map
	replacePlaceholders(configMap)   // Replace placeholders in the map

	// Write back the modified config to Viper
	for key, value := range configMap {
		viper.Set(key, value)
	}

	// Unmarshal into the Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}

// Replace placeholders in a map recursively
func replacePlaceholders(configMap map[string]interface{}) {
	for key, value := range configMap {
		switch v := value.(type) {
		case string:
			if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
				envVar := strings.TrimSuffix(strings.TrimPrefix(v, "${"), "}")
				configMap[key] = getEnv(envVar, v) // Replace with env var value or keep as-is
			}
		case map[string]interface{}:
			replacePlaceholders(v) // Recurse for nested maps
		}
	}
}

// Helper function to get an environment variable value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
