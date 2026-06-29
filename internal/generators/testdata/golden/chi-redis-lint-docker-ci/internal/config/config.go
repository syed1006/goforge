// Package config loads runtime configuration from the environment.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the runtime configuration loaded from the environment.
type Config struct {
	HTTPAddr     string
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// Load reads config from the environment, falling back to dev defaults.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:     getenv("HTTP_ADDR", ":8080"),
		LogLevel:     getenv("LOG_LEVEL", "info"),
		ReadTimeout:  getDuration("HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
		RedisAddr:     getenv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       getInt("REDIS_DB", 0),
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

func getInt(k string, fallback int) int {
	v := os.Getenv(k)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

