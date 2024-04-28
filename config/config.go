package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server     *ServerConfig
	Auth       *AuthConfig
	PostgreSQL *PostgresConfig
}

type ServerConfig struct {
	Host string
	Port string
	Docs bool
}

type AuthConfig struct {
	Username string
	Password string
}

type PostgresConfig struct {
	URL string
}

func GetConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Host: getEnvDefault("HOST", "0.0.0.0"),
			Port: getEnvDefault("PORT", "8080"),
			Docs: getBoolDefault("DOCS", true),
		},
		Auth: &AuthConfig{
			Username: getEnvDefault("ADMIN_USERNAME", ""),
			Password: getEnvDefault("ADMIN_PASSWORD", ""),
		},
		PostgreSQL: &PostgresConfig{
			URL: getEnvDefault("DATABASE_URL", ""),
		},
	}
}

func getEnvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getBoolDefault(key string, fallback bool) bool {
	value := getEnvDefault(key, strconv.FormatBool(fallback))
	if val, err := strconv.ParseBool(value); err == nil {
		return val
	}

	return fallback
}
