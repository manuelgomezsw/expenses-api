package daily_expense_config

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// DailyExpenseConfig represents monthly budget configuration for daily expenses
// Maps to frontend interface: DailyExpensesConfig { id?, monthly_budget, month }
type DailyExpenseConfig struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	MonthlyBudget float64 `gorm:"type:decimal(15,2);not null" json:"monthly_budget"`
	Month         string  `gorm:"size:7;not null;uniqueIndex" json:"month"` // Format: "2024-01"
}

// TableName specifies the table name for GORM
func (DailyExpenseConfig) TableName() string {
	return "daily_expenses_configs"
}

// BeforeCreate hook to validate data before creation
func (dec *DailyExpenseConfig) BeforeCreate(tx *gorm.DB) error {
	return dec.validate()
}

// BeforeUpdate hook to validate data before update
func (dec *DailyExpenseConfig) BeforeUpdate(tx *gorm.DB) error {
	return dec.validate()
}

// validate performs validation
func (dec *DailyExpenseConfig) validate() error {
	// Validate monthly budget is not negative
	if dec.MonthlyBudget < 0 {
		return errors.New("monthly budget cannot be negative")
	}

	// Validate month format (YYYY-MM)
	if len(dec.Month) != 7 {
		return errors.New("month must be in YYYY-MM format")
	}

	// Try to parse the month to ensure it's valid
	_, err := time.Parse("2006-01", dec.Month)
	if err != nil {
		return errors.New("invalid month format, must be YYYY-MM")
	}

	return nil
}

// GetDailyBudget calculates the daily budget based on the monthly budget
func (dec *DailyExpenseConfig) GetDailyBudget() float64 {
	if dec.MonthlyBudget <= 0 {
		return 0
	}

	// Parse month to get the number of days
	date, err := time.Parse("2006-01", dec.Month)
	if err != nil {
		return 0
	}

	// Get the last day of the month
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())
	daysInMonth := lastDay.Day()

	return dec.MonthlyBudget / float64(daysInMonth)
}

// GetRemainingDays calculates remaining days in the month
func (dec *DailyExpenseConfig) GetRemainingDays() int {
	// Parse month
	date, err := time.Parse("2006-01", dec.Month)
	if err != nil {
		return 0
	}

	// Get current date
	now := time.Now()

	// If it's not the current month, return 0
	if now.Format("2006-01") != dec.Month {
		return 0
	}

	// Get the last day of the month
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())

	return lastDay.Day() - now.Day() + 1
}

// IsCurrentMonth checks if this config is for the current month
func (dec *DailyExpenseConfig) IsCurrentMonth() bool {
	currentMonth := time.Now().Format("2006-01")
	return dec.Month == currentMonth
}

// GetCurrentMonth returns the current month in YYYY-MM format
func GetCurrentMonth() string {
	return time.Now().Format("2006-01")
}
