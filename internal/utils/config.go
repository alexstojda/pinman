package utils

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`

	SPAPath          string   `mapstructure:"SPA_PATH"`
	SPACacheDisabled bool     `mapstructure:"SPA_CACHE_DISABLED"`
	ClientOrigins    []string `mapstructure:"CLIENT_ORIGINS"`

	TokenPrivateKey   string        `mapstructure:"TOKEN_PRIVATE_KEY"`
	TokenPublicKey    string        `mapstructure:"TOKEN_PUBLIC_KEY"`
	TokenExpiresAfter time.Duration `mapstructure:"TOKEN_EXPIRES_AFTER"`
	TokenSecretKey    string        `mapstructure:"TOKEN_SECRET_KEY"`

	// Added to support deployment to railway
	RailwayDBUrl     string `mapstructure:"DATABASE_URL"`
	RailwayStaticUrl string `mapstructure:"RAILWAY_STATIC_URL"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetTypeByDefaultValue(true)
	viper.SetConfigType("env")
	if envFile := os.Getenv("ENV_FILE"); envFile != "" {
		viper.SetConfigFile(envFile)
	} else {
		viper.SetConfigFile(".env")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	// Silly bug in viper, it will only read ENV variables if they are defined in a .env file.
	// So we manually tell viper to read all variables defined in the Config type.
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(config, &envKeysMap); err != nil {
		return nil, err
	}
	for k := range *envKeysMap {
		if bindErr := viper.BindEnv(k); bindErr != nil {
			return nil, err
		}
	}

	err = viper.Unmarshal(config)

	if len(config.RailwayStaticUrl) > 0 {
		config.ClientOrigins = append(config.ClientOrigins, config.RailwayStaticUrl)
	}

	return config, nil
}
