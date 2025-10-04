package port

import (
	"expenses-api/internal/domain/daily_expense"
	"expenses-api/internal/domain/daily_expense_config"
	"expenses-api/internal/domain/fixed_expense"
	"expenses-api/internal/domain/pocket"
	"expenses-api/internal/domain/salary"
)

// SalaryRepository defines the interface for salary data operations
// Frontend endpoints: GET/PUT /api/config/income
type SalaryRepository interface {
	GetByMonth(month string) (*salary.Salary, error)
	CreateOrUpdate(s *salary.Salary) error
}

// PocketRepository defines the interface for pocket data operations
// Frontend endpoints: GET/POST/PUT/DELETE /api/config/pockets
type PocketRepository interface {
	GetAll() ([]pocket.Pocket, error)
	GetByID(id uint) (*pocket.Pocket, error)
	GetByName(name string) (*pocket.Pocket, error)
	Create(p *pocket.Pocket) error
	Update(p *pocket.Pocket) error
	Delete(id uint) error
}

// FixedExpenseRepository defines the interface for fixed expense data operations
// Frontend endpoints: GET /api/fixed-expenses/{month}, POST/PUT /api/fixed-expenses, PUT /api/fixed-expenses/{id}/status
type FixedExpenseRepository interface {
	GetByMonth(month string) ([]fixed_expense.FixedExpense, error)
	GetByID(id uint) (*fixed_expense.FixedExpense, error)
	Create(expense *fixed_expense.FixedExpense) error
	Update(expense *fixed_expense.FixedExpense) error
	UpdatePaymentStatus(id uint, isPaid bool, paidDate *string) error
}

// DailyExpenseRepository defines the interface for daily expense data operations
// Frontend endpoints: GET/POST/PUT/DELETE /api/daily-expenses
type DailyExpenseRepository interface {
	GetByMonth(month string) ([]daily_expense.DailyExpense, error)
	GetByID(id uint) (*daily_expense.DailyExpense, error)
	Create(expense *daily_expense.DailyExpense) error
	Update(expense *daily_expense.DailyExpense) error
	Delete(id uint) error
}

// DailyExpenseConfigRepository defines the interface for daily expense config data operations
// Frontend endpoints: GET/PUT /api/config/daily-budget/{month}
type DailyExpenseConfigRepository interface {
	GetByMonth(month string) (*daily_expense_config.DailyExpenseConfig, error)
	CreateOrUpdate(config *daily_expense_config.DailyExpenseConfig) error
}
