package usecase

import (
	"errors"
	"expenses-api/internal/application/port"
	"expenses-api/internal/api/dto"
	"time"
)

// SummaryUseCase handles summary-related business logic
type SummaryUseCase struct {
	salaryRepo            port.SalaryRepository
	fixedExpenseRepo      port.FixedExpenseRepository
	dailyExpenseRepo      port.DailyExpenseRepository
	dailyExpenseConfigRepo port.DailyExpenseConfigRepository
}

// NewSummaryUseCase creates a new summary use case instance
func NewSummaryUseCase(
	salaryRepo port.SalaryRepository,
	fixedExpenseRepo port.FixedExpenseRepository,
	dailyExpenseRepo port.DailyExpenseRepository,
	dailyExpenseConfigRepo port.DailyExpenseConfigRepository,
) *SummaryUseCase {
	return &SummaryUseCase{
		salaryRepo:            salaryRepo,
		fixedExpenseRepo:      fixedExpenseRepo,
		dailyExpenseRepo:      dailyExpenseRepo,
		dailyExpenseConfigRepo: dailyExpenseConfigRepo,
	}
}

// GetMonthlySummary calculates and returns the monthly financial summary
func (uc *SummaryUseCase) GetMonthlySummary(month string) (*dto.MonthlySummaryDTO, error) {
	if month == "" {
		return nil, errors.New("month is required")
	}
	
	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return nil, errors.New("invalid month format, must be YYYY-MM")
	}
	
	// Get salary for the month
	salary, err := uc.salaryRepo.GetByMonth(month)
	var totalIncome float64 = 0
	if err == nil && salary != nil {
		totalIncome = salary.MonthlyAmount
	}
	
	// Get fixed expenses for the month
	fixedExpenses, err := uc.fixedExpenseRepo.GetByMonth(month)
	if err != nil {
		return nil, err
	}
	
	// Calculate fixed expenses totals
	var totalFixedExpenses float64 = 0
	var fixedExpensesPaid int = 0
	var fixedExpensesTotal int = len(fixedExpenses)
	
	for _, expense := range fixedExpenses {
		totalFixedExpenses += expense.Amount
		if expense.IsPaid {
			fixedExpensesPaid++
		}
	}
	
	// Get daily expenses for the month
	dailyExpenses, err := uc.dailyExpenseRepo.GetByMonth(month)
	if err != nil {
		return nil, err
	}
	
	// Calculate daily expenses total
	var totalDailyExpenses float64 = 0
	for _, expense := range dailyExpenses {
		totalDailyExpenses += expense.Amount
	}
	
	// Get daily expense config for the month
	dailyConfig, err := uc.dailyExpenseConfigRepo.GetByMonth(month)
	var dailyBudgetTotal float64 = 0
	if err == nil && dailyConfig != nil {
		dailyBudgetTotal = dailyConfig.MonthlyBudget
	}
	
	// Calculate remaining budget
	remainingBudget := totalIncome - totalFixedExpenses - totalDailyExpenses
	
	summary := &dto.MonthlySummaryDTO{
		Month:              month,
		TotalIncome:        totalIncome,
		TotalFixedExpenses: totalFixedExpenses,
		TotalDailyExpenses: totalDailyExpenses,
		RemainingBudget:    remainingBudget,
		FixedExpensesPaid:  fixedExpensesPaid,
		FixedExpensesTotal: fixedExpensesTotal,
		DailyBudgetUsed:    totalDailyExpenses,
		DailyBudgetTotal:   dailyBudgetTotal,
	}
	
	return summary, nil
}

// GetCurrentMonthlySummary returns summary for the current month
func (uc *SummaryUseCase) GetCurrentMonthlySummary() (*dto.MonthlySummaryDTO, error) {
	currentMonth := time.Now().Format("2006-01")
	return uc.GetMonthlySummary(currentMonth)
}
