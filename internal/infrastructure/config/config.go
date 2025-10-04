package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port string
	Env  string

	// Database
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBDSN      string // For GCP Cloud SQL

	// Security
	JWTSecret string

	// CORS
	CORSAllowedOrigin string

	// Debug
	Debug    bool
	LogLevel string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Try to load .env file (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		// Server configuration
		Port: getEnvOrDefault("PORT", "8080"),
		Env:  getEnvOrDefault("ENV", "development"),

		// Database configuration
		DBHost:     getEnvOrDefault("DB_HOST", "localhost:3306"),
		DBUser:     getEnvOrDefault("DB_USER", "root"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:     getEnvOrDefault("DB_NAME_EXPENSES", "expenses_db"),
		DBDSN:      getEnvOrDefault("DB_DSN", ""), // For GCP

		// Security
		JWTSecret: getEnvOrDefault("JWT_SECRET", "default-secret-change-in-production"),

		// CORS
		CORSAllowedOrigin: getEnvOrDefault("CORS_ALLOWED_ORIGIN", "http://localhost:4200"),

		// Debug
		Debug:    getEnvAsBool("DEBUG", true),
		LogLevel: getEnvOrDefault("LOG_LEVEL", "info"),
	}

	// Validate required configuration
	if err := config.validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	AppConfig = config
	return config
}

// GetDSN returns the appropriate database connection string
func (c *Config) GetDSN() string {
	// If DB_DSN is provided (GCP), use it directly
	if c.DBDSN != "" {
		return c.DBDSN
	}

	// Otherwise, build DSN from individual components (local development)
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBName)
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Env) == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Env) == "development"
}

// validate checks if all required configuration is present
func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	// Database validation
	if c.DBDSN == "" {
		// If no DSN, validate individual components
		if c.DBPassword == "" {
			return fmt.Errorf("DB_PASSWORD is required")
		}
		if c.DBName == "" {
			return fmt.Errorf("DB_NAME_EXPENSES is required")
		}
	}

	// Security validation
	if c.IsProduction() && c.JWTSecret == "default-secret-change-in-production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	return nil
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets environment variable as boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
