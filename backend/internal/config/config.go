package config

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port           string
	GinMode        string
	FrontendOrigin string
	ScanTimeout    time.Duration
}

func Load() Config {
	return Config{
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", gin.ReleaseMode),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		ScanTimeout:    time.Duration(getEnvInt("SCAN_TIMEOUT_SECONDS", 60)) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
