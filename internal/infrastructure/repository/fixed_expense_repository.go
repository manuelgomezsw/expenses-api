package repository

import (
	"expenses-api/internal/domain/fixed_expense"

	"gorm.io/gorm"
)

// FixedExpenseRepository handles fixed expense-related database operations
type FixedExpenseRepository struct {
	*BaseRepository
}

// NewFixedExpenseRepository creates a new fixed expense repository instance
func NewFixedExpenseRepository(db *gorm.DB) *FixedExpenseRepository {
	return &FixedExpenseRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetByMonth retrieves all fixed expenses for a specific month with pocket information
func (r *FixedExpenseRepository) GetByMonth(month string) ([]fixed_expense.FixedExpense, error) {
	var expenses []fixed_expense.FixedExpense
	err := r.db.Preload("Pocket").
		Where("month = ?", month).
		Order("payment_day ASC, concept_name ASC").
		Find(&expenses).Error
	return expenses, err
}

// GetByMonthAndPocket retrieves fixed expenses for a specific month and pocket
func (r *FixedExpenseRepository) GetByMonthAndPocket(month string, pocketID uint) ([]fixed_expense.FixedExpense, error) {
	var expenses []fixed_expense.FixedExpense
	err := r.db.Preload("Pocket").
		Where("month = ? AND pocket_id = ?", month, pocketID).
		Order("payment_day ASC, concept_name ASC").
		Find(&expenses).Error
	return expenses, err
}

// GetByID retrieves a fixed expense by ID with pocket information
func (r *FixedExpenseRepository) GetByID(id uint) (*fixed_expense.FixedExpense, error) {
	var expense fixed_expense.FixedExpense
	err := r.db.Preload("Pocket").First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

// Create creates a new fixed expense
func (r *FixedExpenseRepository) Create(expense *fixed_expense.FixedExpense) error {
	return r.db.Create(expense).Error
}

// Update updates an existing fixed expense
func (r *FixedExpenseRepository) Update(expense *fixed_expense.FixedExpense) error {
	return r.db.Save(expense).Error
}

// Delete deletes a fixed expense by ID
func (r *FixedExpenseRepository) Delete(id uint) error {
	return r.db.Delete(&fixed_expense.FixedExpense{}, id).Error
}

// UpdatePaymentStatus updates the payment status of a fixed expense
func (r *FixedExpenseRepository) UpdatePaymentStatus(id uint, isPaid bool, paidDate *string) error {
	updates := map[string]interface{}{
		"is_paid": isPaid,
	}

	if isPaid && paidDate != nil {
		updates["paid_date"] = *paidDate
	} else if !isPaid {
		updates["paid_date"] = nil
	}

	return r.db.Model(&fixed_expense.FixedExpense{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// GetPaidByMonth retrieves all paid fixed expenses for a specific month
func (r *FixedExpenseRepository) GetPaidByMonth(month string) ([]fixed_expense.FixedExpense, error) {
	var expenses []fixed_expense.FixedExpense
	err := r.db.Preload("Pocket").
		Where("month = ? AND is_paid = ?", month, true).
		Order("paid_date DESC, concept_name ASC").
		Find(&expenses).Error
	return expenses, err
}

// GetUnpaidByMonth retrieves all unpaid fixed expenses for a specific month
func (r *FixedExpenseRepository) GetUnpaidByMonth(month string) ([]fixed_expense.FixedExpense, error) {
	var expenses []fixed_expense.FixedExpense
	err := r.db.Preload("Pocket").
		Where("month = ? AND is_paid = ?", month, false).
		Order("payment_day ASC, concept_name ASC").
		Find(&expenses).Error
	return expenses, err
}

// GetOverdueByMonth retrieves overdue fixed expenses for a specific month
func (r *FixedExpenseRepository) GetOverdueByMonth(month string, currentDay int) ([]fixed_expense.FixedExpense, error) {
	var expenses []fixed_expense.FixedExpense
	err := r.db.Preload("Pocket").
		Where("month = ? AND is_paid = ? AND payment_day < ?", month, false, currentDay).
		Order("payment_day ASC, concept_name ASC").
		Find(&expenses).Error
	return expenses, err
}

// GetSummaryByMonth calculates summary statistics for fixed expenses in a month
func (r *FixedExpenseRepository) GetSummaryByMonth(month string) (*FixedExpenseSummary, error) {
	var summary FixedExpenseSummary

	err := r.db.Model(&fixed_expense.FixedExpense{}).
		Select(`
			COUNT(*) as total_count,
			SUM(amount) as total_amount,
			SUM(CASE WHEN is_paid = true THEN 1 ELSE 0 END) as paid_count,
			SUM(CASE WHEN is_paid = true THEN amount ELSE 0 END) as paid_amount,
			SUM(CASE WHEN is_paid = false THEN 1 ELSE 0 END) as unpaid_count,
			SUM(CASE WHEN is_paid = false THEN amount ELSE 0 END) as unpaid_amount
		`).
		Where("month = ?", month).
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	summary.Month = month
	return &summary, nil
}

// GetMonthsWithExpenses retrieves all months that have fixed expenses
func (r *FixedExpenseRepository) GetMonthsWithExpenses() ([]string, error) {
	var months []string
	err := r.db.Model(&fixed_expense.FixedExpense{}).
		Select("DISTINCT month").
		Order("month DESC").
		Pluck("month", &months).Error
	return months, err
}

// BulkUpdatePaymentStatus updates payment status for multiple expenses
func (r *FixedExpenseRepository) BulkUpdatePaymentStatus(ids []uint, isPaid bool, paidDate *string) error {
	updates := map[string]interface{}{
		"is_paid": isPaid,
	}

	if isPaid && paidDate != nil {
		updates["paid_date"] = *paidDate
	} else if !isPaid {
		updates["paid_date"] = nil
	}

	return r.db.Model(&fixed_expense.FixedExpense{}).
		Where("id IN ?", ids).
		Updates(updates).Error
}

// GetByPocketAndMonths retrieves fixed expenses for a pocket across multiple months
func (r *FixedExpenseRepository) GetByPocketAndMonths(pocketID uint, months []string) ([]fixed_expense.FixedExpense, error) {
	var expenses []fixed_expense.FixedExpense
	err := r.db.Preload("Pocket").
		Where("pocket_id = ? AND month IN ?", pocketID, months).
		Order("month DESC, payment_day ASC").
		Find(&expenses).Error
	return expenses, err
}

// FixedExpenseSummary represents summary statistics for fixed expenses
type FixedExpenseSummary struct {
	Month        string  `json:"month"`
	TotalCount   int     `json:"total_count"`
	TotalAmount  float64 `json:"total_amount"`
	PaidCount    int     `json:"paid_count"`
	PaidAmount   float64 `json:"paid_amount"`
	UnpaidCount  int     `json:"unpaid_count"`
	UnpaidAmount float64 `json:"unpaid_amount"`
}
