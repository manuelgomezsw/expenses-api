package daily_expense

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// DailyExpense represents daily expenses
// Maps to frontend interface: DailyExpense { id?, description, amount, date, created_at? }
type DailyExpense struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Description string    `gorm:"size:500;not null" json:"description"`
	Amount      float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Date        string    `gorm:"size:10;not null;index" json:"date"` // Format: "2024-01-15"
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for GORM
func (DailyExpense) TableName() string {
	return "daily_expenses"
}

// BeforeCreate hook to validate data before creation
func (de *DailyExpense) BeforeCreate(tx *gorm.DB) error {
	return de.validate()
}

// BeforeUpdate hook to validate data before update
func (de *DailyExpense) BeforeUpdate(tx *gorm.DB) error {
	return de.validate()
}

// validate performs validation and data cleaning
func (de *DailyExpense) validate() error {
	// Clean and validate description
	de.Description = strings.TrimSpace(de.Description)
	if de.Description == "" {
		return errors.New("description cannot be empty")
	}

	if len(de.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}

	// Validate amount
	if de.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Validate date format (YYYY-MM-DD)
	if len(de.Date) != 10 {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	// Try to parse the date to ensure it's valid
	_, err := time.Parse("2006-01-02", de.Date)
	if err != nil {
		return errors.New("invalid date format, must be YYYY-MM-DD")
	}

	return nil
}

// GetMonth returns the month of the expense in YYYY-MM format
func (de *DailyExpense) GetMonth() string {
	if len(de.Date) >= 7 {
		return de.Date[:7]
	}
	return ""
}

// GetDayName returns the day name of the expense date
func (de *DailyExpense) GetDayName() string {
	date, err := time.Parse("2006-01-02", de.Date)
	if err != nil {
		return ""
	}
	return date.Weekday().String()
}

// GetDayNumber returns the day number of the expense date
func (de *DailyExpense) GetDayNumber() int {
	date, err := time.Parse("2006-01-02", de.Date)
	if err != nil {
		return 0
	}
	return date.Day()
}

// IsToday checks if the expense is from today
func (de *DailyExpense) IsToday() bool {
	today := time.Now().Format("2006-01-02")
	return de.Date == today
}

// IsThisMonth checks if the expense is from the current month
func (de *DailyExpense) IsThisMonth() bool {
	currentMonth := time.Now().Format("2006-01")
	return de.GetMonth() == currentMonth
}

// GetCurrentDate returns the current date in YYYY-MM-DD format
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// GetCurrentMonth returns the current month in YYYY-MM format
func GetCurrentMonth() string {
	return time.Now().Format("2006-01")
}
