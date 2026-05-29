package infrastructure

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig `mapstructure:"server"`
	Database DBConfig     `mapstructure:"database"`
	Redis    RedisConfig  `mapstructure:"redis"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug or release
}

type DBConfig struct {
	URL string `mapstructure:"url"`
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // Important: allow override with env vars

	// Default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.url", "postgres://fintech:fintech123@localhost:5433/fintech_ledger?sslmode=disable")
	viper.SetDefault("redis.addr", "localhost:6380")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config.yaml found, using environment variables + defaults")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}
