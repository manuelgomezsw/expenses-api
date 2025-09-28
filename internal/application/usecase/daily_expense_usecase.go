package usecase

import (
	"errors"
	"expenses-api/internal/application/port"
	"expenses-api/internal/domain/daily_expense"
	"strings"
	"time"
)

// DailyExpenseUseCase handles daily expense-related business logic
type DailyExpenseUseCase struct {
	dailyExpenseRepo port.DailyExpenseRepository
}

// NewDailyExpenseUseCase creates a new daily expense use case instance
func NewDailyExpenseUseCase(dailyExpenseRepo port.DailyExpenseRepository) *DailyExpenseUseCase {
	return &DailyExpenseUseCase{
		dailyExpenseRepo: dailyExpenseRepo,
	}
}

// GetByMonth retrieves all daily expenses for a specific month
func (uc *DailyExpenseUseCase) GetByMonth(month string) ([]daily_expense.DailyExpense, error) {
	if month == "" {
		return nil, errors.New("month is required")
	}
	
	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return nil, errors.New("invalid month format, must be YYYY-MM")
	}
	
	return uc.dailyExpenseRepo.GetByMonth(month)
}

// GetByID retrieves a daily expense by ID
func (uc *DailyExpenseUseCase) GetByID(id uint) (*daily_expense.DailyExpense, error) {
	if id == 0 {
		return nil, errors.New("expense ID is required")
	}
	
	return uc.dailyExpenseRepo.GetByID(id)
}

// Create creates a new daily expense
func (uc *DailyExpenseUseCase) Create(
	description string,
	amount float64,
	date string,
) (*daily_expense.DailyExpense, error) {
	// Validate input
	description = strings.TrimSpace(description)
	if description == "" {
		return nil, errors.New("description is required")
	}
	
	if len(description) > 500 {
		return nil, errors.New("description cannot exceed 500 characters")
	}
	
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	
	if date == "" {
		return nil, errors.New("date is required")
	}
	
	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, errors.New("invalid date format, must be YYYY-MM-DD")
	}
	
	// Don't allow future dates beyond today
	expenseDate, _ := time.Parse("2006-01-02", date)
	today := time.Now()
	if expenseDate.After(today) {
		return nil, errors.New("expense date cannot be in the future")
	}
	
	// Create daily expense
	expense := &daily_expense.DailyExpense{
		Description: description,
		Amount:      amount,
		Date:        date,
	}
	
	if err := uc.dailyExpenseRepo.Create(expense); err != nil {
		return nil, err
	}
	
	return expense, nil
}

// Update updates an existing daily expense
func (uc *DailyExpenseUseCase) Update(
	id uint,
	description string,
	amount float64,
	date string,
) (*daily_expense.DailyExpense, error) {
	if id == 0 {
		return nil, errors.New("expense ID is required")
	}
	
	// Get existing expense
	existingExpense, err := uc.dailyExpenseRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// Validate input
	description = strings.TrimSpace(description)
	if description == "" {
		return nil, errors.New("description is required")
	}
	
	if len(description) > 500 {
		return nil, errors.New("description cannot exceed 500 characters")
	}
	
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	
	if date == "" {
		return nil, errors.New("date is required")
	}
	
	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, errors.New("invalid date format, must be YYYY-MM-DD")
	}
	
	// Don't allow future dates beyond today
	expenseDate, _ := time.Parse("2006-01-02", date)
	today := time.Now()
	if expenseDate.After(today) {
		return nil, errors.New("expense date cannot be in the future")
	}
	
	// Update expense
	existingExpense.Description = description
	existingExpense.Amount = amount
	existingExpense.Date = date
	
	if err := uc.dailyExpenseRepo.Update(existingExpense); err != nil {
		return nil, err
	}
	
	return existingExpense, nil
}

// Delete deletes a daily expense
func (uc *DailyExpenseUseCase) Delete(id uint) error {
	if id == 0 {
		return errors.New("expense ID is required")
	}
	
	// Verify expense exists
	_, err := uc.dailyExpenseRepo.GetByID(id)
	if err != nil {
		return err
	}
	
	return uc.dailyExpenseRepo.Delete(id)
}