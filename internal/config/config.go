package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
	Env        string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "ecommerce"),
		DBPassword: getEnv("DB_PASSWORD", "ecommerce123"),
		DBName:     getEnv("DB_NAME", "ecommerce"),
		JWTSecret:  getEnv("JWT_SECRET", "balaji-shop-super-secret-jwt-key-change-in-production"),
		Port:       getEnv("PORT", "8080"),
		Env:        getEnv("ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c *Config) DBDSN() string {
	return "host=" + c.DBHost + " user=" + c.DBUser + " password=" + c.DBPassword +
		" dbname=" + c.DBName + " port=" + c.DBPort + " sslmode=disable TimeZone=UTC"
}

func (c *Config) GetPort() int {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8080
	}
	return port
}
