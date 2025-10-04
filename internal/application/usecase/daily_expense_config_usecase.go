package usecase

import (
	"errors"
	"expenses-api/internal/application/port"
	"expenses-api/internal/domain/daily_expense_config"
	"time"
)

// DailyExpenseConfigUseCase handles daily expense config-related business logic
type DailyExpenseConfigUseCase struct {
	dailyExpenseConfigRepo port.DailyExpenseConfigRepository
}

// NewDailyExpenseConfigUseCase creates a new daily expense config use case instance
func NewDailyExpenseConfigUseCase(dailyExpenseConfigRepo port.DailyExpenseConfigRepository) *DailyExpenseConfigUseCase {
	return &DailyExpenseConfigUseCase{
		dailyExpenseConfigRepo: dailyExpenseConfigRepo,
	}
}

// GetByMonth retrieves daily expense configuration for a specific month
func (uc *DailyExpenseConfigUseCase) GetByMonth(month string) (*daily_expense_config.DailyExpenseConfig, error) {
	if month == "" {
		return nil, errors.New("month is required")
	}

	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return nil, errors.New("invalid month format, must be YYYY-MM")
	}

	return uc.dailyExpenseConfigRepo.GetByMonth(month)
}

// UpdateBudget updates or creates the daily expense budget configuration for a specific month
func (uc *DailyExpenseConfigUseCase) UpdateBudget(monthlyBudget float64, month string) error {
	if month == "" {
		return errors.New("month is required")
	}

	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return errors.New("invalid month format, must be YYYY-MM")
	}

	// Validate monthly budget
	if monthlyBudget < 0 {
		return errors.New("monthly budget cannot be negative")
	}

	// Create config object
	config := &daily_expense_config.DailyExpenseConfig{
		MonthlyBudget: monthlyBudget,
		Month:         month,
	}

	return uc.dailyExpenseConfigRepo.CreateOrUpdate(config)
}
