package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBPath     string
	JWTsecret  string
	JWTExpire  int // minutes
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "./data/wenxi.db"),
		JWTsecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpire:  getEnvInt("JWT_EXPIRE_MINUTES", 30),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}