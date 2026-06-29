// Package config loads runtime configuration from the environment.
package config

import (
	"os"
	"time"
)

// Config holds the runtime configuration loaded from the environment.
type Config struct {
	HTTPAddr     string
	GRPCAddr     string
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DatabaseURL  string
}

// Load reads config from the environment, falling back to dev defaults.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:     getenv("HTTP_ADDR", ":8080"),
		GRPCAddr:     getenv("GRPC_ADDR", ":9090"),
		LogLevel:     getenv("LOG_LEVEL", "info"),
		ReadTimeout:  getDuration("HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
	}
	return cfg, nil
}

func getenv(k, fallback string) string {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		return v
	}
	return fallback
}

func getDuration(k string, fallback time.Duration) time.Duration {
	v := os.Getenv(k)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

