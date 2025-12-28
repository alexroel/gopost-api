package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
}

var AppConfig *Config

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	AppConfig = &Config{
		Port:        getEnv("PORT", ":8080"),
		JWTSecret:   getEnv("JWT_SECRET", "default_secret_key"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	return AppConfig
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
