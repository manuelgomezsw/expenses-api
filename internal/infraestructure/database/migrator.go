package database

import (
	"expenses-api/internal/domain/daily_expense"
	"expenses-api/internal/domain/daily_expense_config"
	"expenses-api/internal/domain/fixed_expense"
	"expenses-api/internal/domain/pocket"
	"expenses-api/internal/domain/salary"
	"log"

	"gorm.io/gorm"
)

// Migrator handles database migrations
type Migrator struct {
	db *gorm.DB
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// MigrateAll runs all database migrations
func (m *Migrator) MigrateAll() error {
	log.Println("Starting database migrations...")

	// Define migration order (respecting foreign key dependencies)
	models := []interface{}{
		&salary.Salary{},
		&pocket.Pocket{},
		&fixed_expense.FixedExpense{},
		&daily_expense.DailyExpense{},
		&daily_expense_config.DailyExpenseConfig{},
	}

	// Run auto-migration for all models
	for _, model := range models {
		if err := m.db.AutoMigrate(model); err != nil {
			log.Printf("Failed to migrate %T: %v", model, err)
			return err
		}
		log.Printf("Successfully migrated %T", model)
	}

	log.Println("All migrations completed successfully!")
	return nil
}

// SeedInitialData inserts initial data into the database
func (m *Migrator) SeedInitialData() error {
	log.Println("Starting initial data seeding...")

	// Seed pockets
	if err := m.seedPockets(); err != nil {
		return err
	}

	// Seed current month configurations
	if err := m.seedCurrentMonthConfigs(); err != nil {
		return err
	}

	log.Println("Initial data seeding completed successfully!")
	return nil
}

// seedPockets creates initial pocket categories
func (m *Migrator) seedPockets() error {
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
		err := m.db.Where("name = ?", p.Name).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// Create new pocket
			if err := m.db.Create(&p).Error; err != nil {
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
func (m *Migrator) seedCurrentMonthConfigs() error {
	currentMonth := salary.GetCurrentMonth()

	// Seed salary config
	var existingSalary salary.Salary
	err := m.db.Where("month = ?", currentMonth).First(&existingSalary).Error

	if err == gorm.ErrRecordNotFound {
		salaryConfig := salary.Salary{
			MonthlyAmount: 0.00,
			Month:         currentMonth,
		}
		if err := m.db.Create(&salaryConfig).Error; err != nil {
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
	err = m.db.Where("month = ?", currentMonth).First(&existingConfig).Error

	if err == gorm.ErrRecordNotFound {
		dailyConfig := daily_expense_config.DailyExpenseConfig{
			MonthlyBudget: 0.00,
			Month:         currentMonth,
		}
		if err := m.db.Create(&dailyConfig).Error; err != nil {
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

// DropAllTables drops all tables (use with caution!)
func (m *Migrator) DropAllTables() error {
	log.Println("WARNING: Dropping all tables...")

	models := []interface{}{
		&daily_expense_config.DailyExpenseConfig{},
		&daily_expense.DailyExpense{},
		&fixed_expense.FixedExpense{},
		&pocket.Pocket{},
		&salary.Salary{},
	}

	for _, model := range models {
		if err := m.db.Migrator().DropTable(model); err != nil {
			log.Printf("Failed to drop table for %T: %v", model, err)
			return err
		}
		log.Printf("Dropped table for %T", model)
	}

	log.Println("All tables dropped!")
	return nil
}

// ResetDatabase drops all tables and recreates them with initial data
func (m *Migrator) ResetDatabase() error {
	log.Println("Resetting database...")

	if err := m.DropAllTables(); err != nil {
		return err
	}

	if err := m.MigrateAll(); err != nil {
		return err
	}

	if err := m.SeedInitialData(); err != nil {
		return err
	}

	log.Println("Database reset completed!")
	return nil
}
