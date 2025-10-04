package usecase

import (
	"errors"
	"expenses-api/internal/application/port"
	"expenses-api/internal/domain/salary"
	"strings"
	"time"
)

// SalaryUseCase handles salary-related business logic
type SalaryUseCase struct {
	salaryRepo port.SalaryRepository
}

// NewSalaryUseCase creates a new salary use case instance
func NewSalaryUseCase(salaryRepo port.SalaryRepository) *SalaryUseCase {
	return &SalaryUseCase{
		salaryRepo: salaryRepo,
	}
}

// GetByMonth retrieves salary configuration for a specific month
func (uc *SalaryUseCase) GetByMonth(month string) (*salary.Salary, error) {
	if month == "" {
		return nil, errors.New("month is required")
	}
	
	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return nil, errors.New("invalid month format, must be YYYY-MM")
	}
	
	return uc.salaryRepo.GetByMonth(month)
}

// GetCurrentMonth retrieves salary for the current month
func (uc *SalaryUseCase) GetCurrentMonth() (*salary.Salary, error) {
	currentMonth := salary.GetCurrentMonth()
	return uc.salaryRepo.GetByMonth(currentMonth)
}

// UpdateSalary updates or creates salary configuration for a month
func (uc *SalaryUseCase) UpdateSalary(monthlyAmount float64, month string, currency string) error {
	if month == "" {
		return errors.New("month is required")
	}
	
	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return errors.New("invalid month format, must be YYYY-MM")
	}
	
	if monthlyAmount < 0 {
		return errors.New("monthly amount cannot be negative")
	}
	
	// Set default currency if not provided or empty
	currency = strings.TrimSpace(currency)
	if currency == "" {
		currency = "COP"
	}
	
	// Validate currency format (should be 3 characters)
	if len(currency) != 3 {
		return errors.New("currency must be a 3-character code")
	}
	
	// Convert to uppercase for consistency
	currency = strings.ToUpper(currency)
	
	// Note: Currency is not stored in database, only validated for business logic
	// The domain model only stores MonthlyAmount and Month
	salaryConfig := &salary.Salary{
		MonthlyAmount: monthlyAmount,
		Month:         month,
	}
	
	return uc.salaryRepo.CreateOrUpdate(salaryConfig)
}

