package salary

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Salary represents monthly salary configuration
// Maps to frontend interface: Salary { id?, monthly_amount, month, created_at? }
type Salary struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	MonthlyAmount float64 `gorm:"type:decimal(15,2);not null" json:"monthly_amount"`
	Month         string  `gorm:"size:7;not null;uniqueIndex" json:"month"` // Format: "2024-01"
}

// TableName specifies the table name for GORM
func (Salary) TableName() string {
	return "salaries"
}

// BeforeCreate hook to validate data before creation
func (s *Salary) BeforeCreate(tx *gorm.DB) error {
	// Validate month format (YYYY-MM)
	if len(s.Month) != 7 {
		return errors.New("month must be in YYYY-MM format")
	}

	// Validate monthly amount is not negative
	if s.MonthlyAmount < 0 {
		return errors.New("monthly amount cannot be negative")
	}

	return nil
}

// GetCurrentMonth returns the current month in YYYY-MM format
func GetCurrentMonth() string {
	return time.Now().Format("2006-01")
}
