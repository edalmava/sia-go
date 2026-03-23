package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Env          string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type JWTConfig struct {
	Secret           string
	AccessTTLMinutes int
	RefreshHours     int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() *Config {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	env := getEnv("APP_ENV", "development")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("FATAL: JWT_SECRET environment variable is required for security")
	}

	return &Config{
		Env: env,
		Server: ServerConfig{
			Env:          env,
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		JWT: JWTConfig{
			Secret:           jwtSecret,
			AccessTTLMinutes: getEnvInt("JWT_ACCESS_TTL_MINUTES", 15),
			RefreshHours:     getEnvInt("JWT_REFRESH_HOURS", 168),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "sia"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
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
