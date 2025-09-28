package repository

import (
	"expenses-api/internal/domain/daily_expense_config"

	"gorm.io/gorm"
)

// DailyExpenseConfigRepository handles daily expense config-related database operations
type DailyExpenseConfigRepository struct {
	*BaseRepository
}

// NewDailyExpenseConfigRepository creates a new daily expense config repository instance
func NewDailyExpenseConfigRepository(db *gorm.DB) *DailyExpenseConfigRepository {
	return &DailyExpenseConfigRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetByMonth retrieves daily expense configuration for a specific month
func (r *DailyExpenseConfigRepository) GetByMonth(month string) (*daily_expense_config.DailyExpenseConfig, error) {
	var config daily_expense_config.DailyExpenseConfig
	err := r.db.Where("month = ?", month).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// CreateOrUpdate creates a new config record or updates existing one
func (r *DailyExpenseConfigRepository) CreateOrUpdate(config *daily_expense_config.DailyExpenseConfig) error {
	// Try to find existing record
	var existing daily_expense_config.DailyExpenseConfig
	err := r.db.Where("month = ?", config.Month).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new record
		return r.db.Create(config).Error
	} else if err != nil {
		// Other error
		return err
	}

	// Update existing record
	existing.MonthlyBudget = config.MonthlyBudget
	return r.db.Save(&existing).Error
}

// GetAll retrieves all daily expense configurations ordered by month descending
func (r *DailyExpenseConfigRepository) GetAll() ([]daily_expense_config.DailyExpenseConfig, error) {
	var configs []daily_expense_config.DailyExpenseConfig
	err := r.db.Order("month DESC").Find(&configs).Error
	return configs, err
}

// GetRecent retrieves the most recent config records (limit specified)
func (r *DailyExpenseConfigRepository) GetRecent(limit int) ([]daily_expense_config.DailyExpenseConfig, error) {
	var configs []daily_expense_config.DailyExpenseConfig
	err := r.db.Order("month DESC").Limit(limit).Find(&configs).Error
	return configs, err
}

// DeleteByMonth deletes config record for a specific month
func (r *DailyExpenseConfigRepository) DeleteByMonth(month string) error {
	return r.db.Where("month = ?", month).Delete(&daily_expense_config.DailyExpenseConfig{}).Error
}

// GetCurrentMonth retrieves config for the current month
func (r *DailyExpenseConfigRepository) GetCurrentMonth() (*daily_expense_config.DailyExpenseConfig, error) {
	currentMonth := daily_expense_config.GetCurrentMonth()
	return r.GetByMonth(currentMonth)
}

// GetMonthsWithConfig retrieves all months that have budget configured
func (r *DailyExpenseConfigRepository) GetMonthsWithConfig() ([]string, error) {
	var months []string
	err := r.db.Model(&daily_expense_config.DailyExpenseConfig{}).
		Select("month").
		Order("month DESC").
		Pluck("month", &months).Error
	return months, err
}

// GetTotalBudgetByMonths calculates total budget for multiple months
func (r *DailyExpenseConfigRepository) GetTotalBudgetByMonths(months []string) (float64, error) {
	var total float64
	err := r.db.Model(&daily_expense_config.DailyExpenseConfig{}).
		Select("COALESCE(SUM(monthly_budget), 0)").
		Where("month IN ?", months).
		Scan(&total).Error
	return total, err
}

// GetConfigsWithUsage retrieves configs with actual usage statistics
func (r *DailyExpenseConfigRepository) GetConfigsWithUsage() ([]ConfigWithUsage, error) {
	var results []ConfigWithUsage

	err := r.db.Table("daily_expenses_configs dec").
		Select(`
			dec.id,
			dec.monthly_budget,
			dec.month,
			dec.created_at,
			COALESCE(de_stats.total_spent, 0) as total_spent,
			COALESCE(de_stats.expense_count, 0) as expense_count,
			CASE 
				WHEN dec.monthly_budget > 0 THEN 
					ROUND((COALESCE(de_stats.total_spent, 0) / dec.monthly_budget) * 100, 2)
				ELSE 0
			END as usage_percentage
		`).
		Joins(`
			LEFT JOIN (
				SELECT 
					LEFT(date, 7) as month,
					SUM(amount) as total_spent,
					COUNT(*) as expense_count
				FROM daily_expenses
				GROUP BY LEFT(date, 7)
			) de_stats ON dec.month = de_stats.month
		`).
		Order("dec.month DESC").
		Scan(&results).Error

	return results, err
}

// GetBudgetUtilization calculates budget utilization for a specific month
func (r *DailyExpenseConfigRepository) GetBudgetUtilization(month string) (*BudgetUtilization, error) {
	var result BudgetUtilization

	err := r.db.Table("daily_expenses_configs dec").
		Select(`
			dec.monthly_budget,
			COALESCE(de_stats.total_spent, 0) as total_spent,
			COALESCE(de_stats.expense_count, 0) as expense_count,
			(dec.monthly_budget - COALESCE(de_stats.total_spent, 0)) as remaining_budget,
			CASE 
				WHEN dec.monthly_budget > 0 THEN 
					ROUND((COALESCE(de_stats.total_spent, 0) / dec.monthly_budget) * 100, 2)
				ELSE 0
			END as usage_percentage
		`).
		Joins(`
			LEFT JOIN (
				SELECT 
					LEFT(date, 7) as month,
					SUM(amount) as total_spent,
					COUNT(*) as expense_count
				FROM daily_expenses
				WHERE LEFT(date, 7) = ?
				GROUP BY LEFT(date, 7)
			) de_stats ON dec.month = de_stats.month
		`, month).
		Where("dec.month = ?", month).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	result.Month = month
	return &result, nil
}

// ConfigWithUsage represents a config with usage statistics
type ConfigWithUsage struct {
	ID              uint    `json:"id"`
	MonthlyBudget   float64 `json:"monthly_budget"`
	Month           string  `json:"month"`
	CreatedAt       string  `json:"created_at"`
	TotalSpent      float64 `json:"total_spent"`
	ExpenseCount    int     `json:"expense_count"`
	UsagePercentage float64 `json:"usage_percentage"`
}

// BudgetUtilization represents budget utilization statistics
type BudgetUtilization struct {
	Month           string  `json:"month"`
	MonthlyBudget   float64 `json:"monthly_budget"`
	TotalSpent      float64 `json:"total_spent"`
	ExpenseCount    int     `json:"expense_count"`
	RemainingBudget float64 `json:"remaining_budget"`
	UsagePercentage float64 `json:"usage_percentage"`
}
