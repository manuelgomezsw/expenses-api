package usecase

import (
	"errors"
	"expenses-api/internal/application/port"
	"expenses-api/internal/domain/fixed_expense"
	"time"
)

// FixedExpenseUseCase handles fixed expense-related business logic
type FixedExpenseUseCase struct {
	fixedExpenseRepo port.FixedExpenseRepository
}

// NewFixedExpenseUseCase creates a new fixed expense use case instance
func NewFixedExpenseUseCase(fixedExpenseRepo port.FixedExpenseRepository) *FixedExpenseUseCase {
	return &FixedExpenseUseCase{
		fixedExpenseRepo: fixedExpenseRepo,
	}
}

// GetByMonth retrieves all fixed expenses for a specific month
func (uc *FixedExpenseUseCase) GetByMonth(month string) ([]fixed_expense.FixedExpense, error) {
	if month == "" {
		return nil, errors.New("month is required")
	}
	
	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return nil, errors.New("invalid month format, must be YYYY-MM")
	}
	
	return uc.fixedExpenseRepo.GetByMonth(month)
}

// UpdatePaymentStatus updates the payment status of a fixed expense
func (uc *FixedExpenseUseCase) UpdatePaymentStatus(id uint, isPaid bool) error {
	if id == 0 {
		return errors.New("expense ID is required")
	}
	
	var paidDate *string
	if isPaid {
		currentDate := time.Now().Format("2006-01-02")
		paidDate = &currentDate
	}
	
	return uc.fixedExpenseRepo.UpdatePaymentStatus(id, isPaid, paidDate)
}