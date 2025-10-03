package usecase

import (
	"errors"
	"expenses-api/internal/application/port"
	"expenses-api/internal/domain/fixed_expense"
	"fmt"
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

// Create creates a new fixed expense
func (uc *FixedExpenseUseCase) Create(expense *fixed_expense.FixedExpense) error {
	if expense == nil {
		return errors.New("expense is required")
	}

	// Validate required fields
	if expense.ConceptName == "" {
		return errors.New("concept name is required")
	}
	if expense.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if expense.PaymentDay < 1 || expense.PaymentDay > 31 {
		return errors.New("payment day must be between 1 and 31")
	}
	if expense.Month == "" {
		return errors.New("month is required")
	}
	if expense.PocketID == 0 {
		return errors.New("pocket ID is required")
	}

	// Validate month format
	if _, err := time.Parse("2006-01", expense.Month); err != nil {
		return errors.New("invalid month format, must be YYYY-MM")
	}

	// Set default values
	expense.PaidDate = nil

	return uc.fixedExpenseRepo.Create(expense)
}

// Update updates an existing fixed expense
func (uc *FixedExpenseUseCase) Update(id uint, updatedExpense *fixed_expense.FixedExpense) error {
	if id == 0 {
		return errors.New("expense ID is required")
	}
	if updatedExpense == nil {
		return errors.New("expense data is required")
	}

	// Get existing expense
	existingExpense, err := uc.fixedExpenseRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("expense not found: %w", err)
	}

	// Validate required fields
	if updatedExpense.ConceptName == "" {
		return errors.New("concept name is required")
	}
	if updatedExpense.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if updatedExpense.PaymentDay < 1 || updatedExpense.PaymentDay > 31 {
		return errors.New("payment day must be between 1 and 31")
	}
	if updatedExpense.Month == "" {
		return errors.New("month is required")
	}
	if updatedExpense.PocketID == 0 {
		return errors.New("pocket ID is required")
	}

	// Validate month format
	if _, err := time.Parse("2006-01", updatedExpense.Month); err != nil {
		return errors.New("invalid month format, must be YYYY-MM")
	}

	// Update fields
	existingExpense.ConceptName = updatedExpense.ConceptName
	existingExpense.Amount = updatedExpense.Amount
	existingExpense.PaymentDay = updatedExpense.PaymentDay
	existingExpense.Month = updatedExpense.Month
	existingExpense.PocketID = updatedExpense.PocketID

	// Don't update payment status through this method
	// Use UpdatePaymentStatus for that

	return uc.fixedExpenseRepo.Update(existingExpense)
}

// GetByID retrieves a fixed expense by ID
func (uc *FixedExpenseUseCase) GetByID(id uint) (*fixed_expense.FixedExpense, error) {
	if id == 0 {
		return nil, errors.New("expense ID is required")
	}

	return uc.fixedExpenseRepo.GetByID(id)
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
