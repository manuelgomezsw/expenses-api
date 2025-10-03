package database

import (
	"expenses-api/internal/infrastructure/config"
	"fmt"
	"log"
	"os"
	"strings"
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
	cfg := config.AppConfig
	if cfg == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	dsn := cfg.GetDSN()
	log.Printf("Connecting to database with DSN: %s", maskPassword(dsn))

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

// maskPassword masks the password in DSN for logging
func maskPassword(dsn string) string {
	// Simple password masking for logs
	// Format: user:password@tcp(host)/db
	if idx := strings.Index(dsn, ":"); idx != -1 {
		if idx2 := strings.Index(dsn[idx:], "@"); idx2 != -1 {
			return dsn[:idx+1] + "****" + dsn[idx+idx2:]
		}
	}
	return dsn
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
