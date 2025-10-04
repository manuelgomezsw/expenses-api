package dto

import "time"

// DTOs específicos para el frontend Angular

// SalaryDTO representa la configuración de salario para el frontend
type SalaryDTO struct {
	MonthlyAmount float64 `json:"monthly_amount" binding:"required,min=0"`
}

// FixedExpenseDTO representa un gasto fijo para el frontend
type FixedExpenseDTO struct {
	ID          int     `json:"id"`
	PocketName  string  `json:"pocket_name"`
	ConceptName string  `json:"concept_name" binding:"required,min=1,max=255"`
	Amount      float64 `json:"amount" binding:"required,min=0"`
	PaymentDay  int     `json:"payment_day" binding:"required,min=1,max=31"`
	Month       string  `json:"month"`
	IsPaid      bool    `json:"is_paid"`
	PaidDate    *string `json:"paid_date"`
	PocketID    int     `json:"pocket_id,omitempty" binding:"omitempty,min=1"` // Solo para operaciones de escritura
}

// DailyExpenseDTO representa un gasto diario para el frontend
type DailyExpenseDTO struct {
	ID          int       `json:"id"`
	Amount      float64   `json:"amount" binding:"required,min=0"`
	Description string    `json:"description" binding:"required,min=1,max=255"`
	Date        string    `json:"date,omitempty"`       // Opcional, se asigna automáticamente a la fecha actual
	CreatedAt   time.Time `json:"created_at,omitempty"` // Timestamp de creación
}

// PocketDTO representa un bolsillo para el frontend
type PocketDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"required,min=1"`
}

// DailyExpensesConfigDTO representa la configuración de gastos diarios
type DailyExpensesConfigDTO struct {
	MonthlyBudget float64 `json:"monthly_budget" binding:"required,min=0"`
}

// MonthlySummaryDTO representa el resumen mensual para el dashboard
type MonthlySummaryDTO struct {
	Month              string  `json:"month"`
	TotalIncome        float64 `json:"total_income"`
	TotalFixedExpenses float64 `json:"total_fixed_expenses"`
	TotalDailyExpenses float64 `json:"total_daily_expenses"`
	RemainingBudget    float64 `json:"remaining_budget"`
	FixedExpensesPaid  int     `json:"fixed_expenses_paid"`
	FixedExpensesTotal int     `json:"fixed_expenses_total"`
	DailyBudgetUsed    float64 `json:"daily_budget_used"`
	DailyBudgetTotal   float64 `json:"daily_budget_total"`
}

// ExpenseStatus representa los posibles estados de un gasto fijo
type ExpenseStatus string

const (
	ExpenseStatusPending ExpenseStatus = "pending"
	ExpenseStatusPaid    ExpenseStatus = "paid"
	ExpenseStatusOverdue ExpenseStatus = "overdue"
)
