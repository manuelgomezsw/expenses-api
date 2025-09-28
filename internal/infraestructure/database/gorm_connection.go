package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB   *gorm.DB
	once sync.Once
)

// Database wraps the GORM database connection
type Database struct {
	DB *gorm.DB
}

// GetDB returns the singleton GORM database instance
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		DB, err = initializeGORM()
		if err != nil {
			log.Fatalf("Failed to initialize GORM: %v", err)
		}
	})
	return DB
}

// initializeGORM creates and configures the GORM database connection
func initializeGORM() (*gorm.DB, error) {
	dsn, err := buildDSN()
	if err != nil {
		return nil, fmt.Errorf("failed to build DSN: %w", err)
	}

	// Configure GORM
	config := &gorm.Config{
		Logger: getLoggerConfig(),
		NowFunc: func() time.Time {
			// Use Colombia timezone
			loc, _ := time.LoadLocation("America/Bogota")
			return time.Now().In(loc)
		},
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("GORM database connection established successfully")
	return db, nil
}

// buildDSN constructs the database connection string from environment variables
func buildDSN() (string, error) {
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "")
	dbHost := getEnvOrDefault("DB_HOST", "localhost:3306")
	dbName := getEnvOrDefault("DB_NAME", "expenses_db")

	if dbPassword == "" {
		return "", fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	// Build DSN for local MySQL or Cloud SQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbName)

	return dsn, nil
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getLoggerConfig returns the appropriate logger configuration based on environment
func getLoggerConfig() logger.Interface {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "development" {
		return logger.Default.LogMode(logger.Info)
	}

	// Production: only log errors
	return logger.Default.LogMode(logger.Error)
}

// NewDatabase creates a new Database instance
func NewDatabase() (*Database, error) {
	db := GetDB()
	return &Database{DB: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Note: Database migrations are handled manually via SQL scripts
// See sql/database/ directory for schema management

// Transaction executes a function within a database transaction
func (d *Database) Transaction(fn func(*gorm.DB) error) error {
	return d.DB.Transaction(fn)
}
