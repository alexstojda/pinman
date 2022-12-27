package utils

import (
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
}

func LoadConfig() (*Config, error) {
	viper.SetTypeByDefaultValue(true)
	viper.SetConfigType("env")
	if envFile := os.Getenv("ENV_FILE"); envFile != "" {
		viper.SetConfigFile(envFile)
	} else {
		viper.SetConfigFile(".env")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	return config, nil
}
