package database

import (
	"expenses-api/internal/domain/daily_expense_config"
	"expenses-api/internal/domain/pocket"
	"expenses-api/internal/domain/salary"
	"log"
	
	"gorm.io/gorm"
)

// Seeder handles initial data seeding (not schema migrations)
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// CheckConnection verifies database connection
func (s *Seeder) CheckConnection() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	
	if err := sqlDB.Ping(); err != nil {
		return err
	}
	
	log.Println("Database connection verified successfully!")
	return nil
}

// SeedInitialData inserts initial data into the database
func (s *Seeder) SeedInitialData() error {
	log.Println("Starting initial data seeding...")

	// Seed pockets
	if err := s.seedPockets(); err != nil {
		return err
	}
	
	// Seed current month configurations
	if err := s.seedCurrentMonthConfigs(); err != nil {
		return err
	}

	log.Println("Initial data seeding completed successfully!")
	return nil
}

// seedPockets creates initial pocket categories
func (s *Seeder) seedPockets() error {
	pockets := []pocket.Pocket{
		{Name: "Hogar", Description: "Gastos relacionados con el hogar y servicios básicos"},
		{Name: "Alimentación", Description: "Comida, supermercado y restaurantes"},
		{Name: "Transporte", Description: "Transporte público, gasolina, mantenimiento vehículo"},
		{Name: "Salud", Description: "Medicina, consultas médicas, seguros de salud"},
		{Name: "Entretenimiento", Description: "Cine, streaming, salidas, hobbies"},
		{Name: "Educación", Description: "Cursos, libros, capacitaciones"},
		{Name: "Ropa", Description: "Vestimenta y accesorios"},
		{Name: "Otros", Description: "Gastos varios no categorizados"},
	}

	for _, p := range pockets {
		// Check if pocket already exists
		var existing pocket.Pocket
		err := s.db.Where("name = ?", p.Name).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// Create new pocket
			if err := s.db.Create(&p).Error; err != nil {
				log.Printf("Failed to create pocket %s: %v", p.Name, err)
				return err
			}
			log.Printf("Created pocket: %s", p.Name)
		} else if err != nil {
			log.Printf("Error checking pocket %s: %v", p.Name, err)
			return err
		}
		// If no error, pocket already exists, skip
	}
	
	return nil
}

// seedCurrentMonthConfigs creates default configurations for current month
func (s *Seeder) seedCurrentMonthConfigs() error {
	currentMonth := salary.GetCurrentMonth()

	// Seed salary config
	var existingSalary salary.Salary
	err := s.db.Where("month = ?", currentMonth).First(&existingSalary).Error
	
	if err == gorm.ErrRecordNotFound {
		salaryConfig := salary.Salary{
			MonthlyAmount: 0.00,
			Month:         currentMonth,
		}
		if err := s.db.Create(&salaryConfig).Error; err != nil {
			log.Printf("Failed to create salary config: %v", err)
			return err
		}
		log.Printf("Created salary config for month: %s", currentMonth)
	} else if err != nil {
		log.Printf("Error checking salary config: %v", err)
		return err
	}
	
	// Seed daily expense config
	var existingConfig daily_expense_config.DailyExpenseConfig
	err = s.db.Where("month = ?", currentMonth).First(&existingConfig).Error
	
	if err == gorm.ErrRecordNotFound {
		dailyConfig := daily_expense_config.DailyExpenseConfig{
			MonthlyBudget: 0.00,
			Month:         currentMonth,
		}
		if err := s.db.Create(&dailyConfig).Error; err != nil {
			log.Printf("Failed to create daily expense config: %v", err)
			return err
		}
		log.Printf("Created daily expense config for month: %s", currentMonth)
	} else if err != nil {
		log.Printf("Error checking daily expense config: %v", err)
		return err
	}

	return nil
}

// Note: Database schema management is handled manually via SQL scripts
// See sql/database/ directory for schema creation and migration scripts
