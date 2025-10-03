package fixed_expense

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// FixedExpense represents monthly fixed expenses
// Maps to frontend interface: FixedExpense { id?, pocket_name, concept_name, amount, payment_day, is_paid, month, paid_date?, created_at? }
type FixedExpense struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	PocketID    uint    `gorm:"not null;index:idx_pocket_month,priority:1" json:"pocket_id"`
	ConceptName string  `gorm:"size:255;not null" json:"concept_name"`
	Amount      float64 `gorm:"type:decimal(15,2);not null" json:"amount"`
	PaymentDay  int     `gorm:"not null;check:payment_day >= 1 AND payment_day <= 31" json:"payment_day"`
	IsPaid      bool    `gorm:"default:false;index" json:"is_paid"`
	Month       string  `gorm:"size:7;not null;index:idx_pocket_month,priority:2" json:"month"` // Format: "2024-01"
	PaidDate    *string `gorm:"size:10" json:"paid_date"`                                       // Format: "2024-01-15"

	// Relationship - will be loaded when needed
	Pocket *Pocket `gorm:"foreignKey:PocketID" json:"pocket,omitempty"`
}

// Pocket represents the relationship to avoid circular imports
type Pocket struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TableName specifies the table name for GORM
func (FixedExpense) TableName() string {
	return "fixed_expenses"
}

// BeforeCreate hook to validate data before creation
func (fe *FixedExpense) BeforeCreate(tx *gorm.DB) error {
	return fe.validate()
}

// BeforeUpdate hook to validate data before update
func (fe *FixedExpense) BeforeUpdate(tx *gorm.DB) error {
	return fe.validate()
}

// validate performs validation and data cleaning
func (fe *FixedExpense) validate() error {
	// Clean and validate concept name
	fe.ConceptName = strings.TrimSpace(fe.ConceptName)
	if fe.ConceptName == "" {
		return errors.New("concept name cannot be empty")
	}

	if len(fe.ConceptName) > 255 {
		return errors.New("concept name cannot exceed 255 characters")
	}

	// Validate amount
	if fe.Amount < 0 {
		return errors.New("amount cannot be negative")
	}

	// Validate payment day
	if fe.PaymentDay < 1 || fe.PaymentDay > 31 {
		return errors.New("payment day must be between 1 and 31")
	}

	// Validate month format (YYYY-MM)
	if len(fe.Month) != 7 {
		return errors.New("month must be in YYYY-MM format")
	}

	// Validate paid date format if provided
	if fe.PaidDate != nil && len(*fe.PaidDate) != 10 {
		return errors.New("paid date must be in YYYY-MM-DD format")
	}

	return nil
}

// MarkAsPaid marks the expense as paid with the current date
func (fe *FixedExpense) MarkAsPaid() {
	fe.IsPaid = true
	currentDate := time.Now().Format("2006-01-02")
	fe.PaidDate = &currentDate
}

// MarkAsUnpaid marks the expense as unpaid
func (fe *FixedExpense) MarkAsUnpaid() {
	fe.IsPaid = false
	fe.PaidDate = nil
}

// GetStatus returns the current status of the fixed expense
func (fe *FixedExpense) GetStatus() string {
	if fe.IsPaid {
		return "paid"
	}

	// Check if overdue (only for current month)
	currentMonth := time.Now().Format("2006-01")
	if fe.Month == currentMonth && time.Now().Day() > fe.PaymentDay {
		return "overdue"
	}

	return "pending"
}

// GetCurrentMonth returns the current month in YYYY-MM format
func GetCurrentMonth() string {
	return time.Now().Format("2006-01")
}
