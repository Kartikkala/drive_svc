package config

import (
	"os"
)

func getEnvOrDefault(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

type DatabaseConfig struct {
	User     string
	Password string
	Hostname string
	DBName   string
	Port     uint16
	Timezone string
}

type NATSConfig struct {
	URL string
}

type Config struct {
	Database DatabaseConfig
	NATS     NATSConfig
}

func NewConfig() *Config {

	return &Config{
		Database: DatabaseConfig{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASS"),
			DBName:   os.Getenv("POSTGRES_DBNAME"),
			Hostname: getEnvOrDefault("POSTGRES_HOSTNAME", "127.0.0.1"),
			Port:     uint16(5432),
			Timezone: getEnvOrDefault("POSTGRES_TIMEZONE", "Asia/Kolkata"),
		},
		NATS: NATSConfig{
			URL: getEnvOrDefault("NATS_URL", "nats://127.0.0.1:4222"),
		},
	}
}
