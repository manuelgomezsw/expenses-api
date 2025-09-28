package repository

import (
	"expenses-api/internal/domain/daily_expense"

	"gorm.io/gorm"
)

// DailyExpenseRepository handles daily expense-related database operations
type DailyExpenseRepository struct {
	*BaseRepository
}

// NewDailyExpenseRepository creates a new daily expense repository instance
func NewDailyExpenseRepository(db *gorm.DB) *DailyExpenseRepository {
	return &DailyExpenseRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetByMonth retrieves all daily expenses for a specific month
func (r *DailyExpenseRepository) GetByMonth(month string) ([]daily_expense.DailyExpense, error) {
	var expenses []daily_expense.DailyExpense
	err := r.db.Where("date LIKE ?", month+"%").
		Order("date DESC, created_at DESC").
		Find(&expenses).Error
	return expenses, err
}

// GetByDateRange retrieves daily expenses within a date range
func (r *DailyExpenseRepository) GetByDateRange(startDate, endDate string) ([]daily_expense.DailyExpense, error) {
	var expenses []daily_expense.DailyExpense
	err := r.db.Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date DESC, created_at DESC").
		Find(&expenses).Error
	return expenses, err
}

// GetByDate retrieves all daily expenses for a specific date
func (r *DailyExpenseRepository) GetByDate(date string) ([]daily_expense.DailyExpense, error) {
	var expenses []daily_expense.DailyExpense
	err := r.db.Where("date = ?", date).
		Order("created_at DESC").
		Find(&expenses).Error
	return expenses, err
}

// GetByID retrieves a daily expense by ID
func (r *DailyExpenseRepository) GetByID(id uint) (*daily_expense.DailyExpense, error) {
	var expense daily_expense.DailyExpense
	err := r.db.First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

// Create creates a new daily expense
func (r *DailyExpenseRepository) Create(expense *daily_expense.DailyExpense) error {
	return r.db.Create(expense).Error
}

// Update updates an existing daily expense
func (r *DailyExpenseRepository) Update(expense *daily_expense.DailyExpense) error {
	return r.db.Save(expense).Error
}

// Delete deletes a daily expense by ID
func (r *DailyExpenseRepository) Delete(id uint) error {
	return r.db.Delete(&daily_expense.DailyExpense{}, id).Error
}

// GetRecent retrieves the most recent daily expenses (limit specified)
func (r *DailyExpenseRepository) GetRecent(limit int) ([]daily_expense.DailyExpense, error) {
	var expenses []daily_expense.DailyExpense
	err := r.db.Order("date DESC, created_at DESC").
		Limit(limit).
		Find(&expenses).Error
	return expenses, err
}

// GetSummaryByMonth calculates summary statistics for daily expenses in a month
func (r *DailyExpenseRepository) GetSummaryByMonth(month string) (*DailyExpenseSummary, error) {
	var summary DailyExpenseSummary

	err := r.db.Model(&daily_expense.DailyExpense{}).
		Select(`
			COUNT(*) as total_count,
			SUM(amount) as total_amount,
			AVG(amount) as average_amount,
			MIN(amount) as min_amount,
			MAX(amount) as max_amount
		`).
		Where("date LIKE ?", month+"%").
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	summary.Month = month
	return &summary, nil
}

// GetDailyTotals retrieves daily totals for a specific month
func (r *DailyExpenseRepository) GetDailyTotals(month string) ([]DailyTotal, error) {
	var totals []DailyTotal

	err := r.db.Model(&daily_expense.DailyExpense{}).
		Select("date, SUM(amount) as total_amount, COUNT(*) as expense_count").
		Where("date LIKE ?", month+"%").
		Group("date").
		Order("date ASC").
		Scan(&totals).Error

	return totals, err
}

// GetMonthsWithExpenses retrieves all months that have daily expenses
func (r *DailyExpenseRepository) GetMonthsWithExpenses() ([]string, error) {
	var months []string
	err := r.db.Model(&daily_expense.DailyExpense{}).
		Select("DISTINCT LEFT(date, 7) as month").
		Order("month DESC").
		Pluck("month", &months).Error
	return months, err
}

// SearchByDescription searches daily expenses by description
func (r *DailyExpenseRepository) SearchByDescription(query string, limit int) ([]daily_expense.DailyExpense, error) {
	var expenses []daily_expense.DailyExpense
	err := r.db.Where("description LIKE ?", "%"+query+"%").
		Order("date DESC, created_at DESC").
		Limit(limit).
		Find(&expenses).Error
	return expenses, err
}

// GetByAmountRange retrieves daily expenses within an amount range
func (r *DailyExpenseRepository) GetByAmountRange(minAmount, maxAmount float64, month string) ([]daily_expense.DailyExpense, error) {
	var expenses []daily_expense.DailyExpense
	query := r.db.Where("amount >= ? AND amount <= ?", minAmount, maxAmount)

	if month != "" {
		query = query.Where("date LIKE ?", month+"%")
	}

	err := query.Order("date DESC, amount DESC").Find(&expenses).Error
	return expenses, err
}

// GetTopExpensesByMonth retrieves the highest expenses for a month
func (r *DailyExpenseRepository) GetTopExpensesByMonth(month string, limit int) ([]daily_expense.DailyExpense, error) {
	var expenses []daily_expense.DailyExpense
	err := r.db.Where("date LIKE ?", month+"%").
		Order("amount DESC, date DESC").
		Limit(limit).
		Find(&expenses).Error
	return expenses, err
}

// BulkDelete deletes multiple daily expenses by IDs
func (r *DailyExpenseRepository) BulkDelete(ids []uint) error {
	return r.db.Delete(&daily_expense.DailyExpense{}, ids).Error
}

// GetExpensesByWeekday retrieves expenses grouped by weekday for a month
func (r *DailyExpenseRepository) GetExpensesByWeekday(month string) ([]WeekdayExpense, error) {
	var results []WeekdayExpense

	err := r.db.Model(&daily_expense.DailyExpense{}).
		Select(`
			DAYNAME(STR_TO_DATE(date, '%Y-%m-%d')) as weekday,
			COUNT(*) as expense_count,
			SUM(amount) as total_amount,
			AVG(amount) as average_amount
		`).
		Where("date LIKE ?", month+"%").
		Group("DAYNAME(STR_TO_DATE(date, '%Y-%m-%d'))").
		Order("FIELD(DAYNAME(STR_TO_DATE(date, '%Y-%m-%d')), 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday')").
		Scan(&results).Error

	return results, err
}

// DailyExpenseSummary represents summary statistics for daily expenses
type DailyExpenseSummary struct {
	Month         string  `json:"month"`
	TotalCount    int     `json:"total_count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
	MinAmount     float64 `json:"min_amount"`
	MaxAmount     float64 `json:"max_amount"`
}

// DailyTotal represents daily expense totals
type DailyTotal struct {
	Date         string  `json:"date"`
	TotalAmount  float64 `json:"total_amount"`
	ExpenseCount int     `json:"expense_count"`
}

// WeekdayExpense represents expenses grouped by weekday
type WeekdayExpense struct {
	Weekday       string  `json:"weekday"`
	ExpenseCount  int     `json:"expense_count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}
