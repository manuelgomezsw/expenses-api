package repository

import (
	"expenses-api/internal/domain/salary"

	"gorm.io/gorm"
)

// SalaryRepository handles salary-related database operations
type SalaryRepository struct {
	*BaseRepository
}

// NewSalaryRepository creates a new salary repository instance
func NewSalaryRepository(db *gorm.DB) *SalaryRepository {
	return &SalaryRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetByMonth retrieves salary configuration for a specific month
func (r *SalaryRepository) GetByMonth(month string) (*salary.Salary, error) {
	var s salary.Salary
	err := r.db.Where("month = ?", month).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateOrUpdate creates a new salary record or updates existing one
func (r *SalaryRepository) CreateOrUpdate(s *salary.Salary) error {
	// Try to find existing record
	var existing salary.Salary
	err := r.db.Where("month = ?", s.Month).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new record
		return r.db.Create(s).Error
	} else if err != nil {
		// Other error
		return err
	}

	// Update existing record
	existing.MonthlyAmount = s.MonthlyAmount
	return r.db.Save(&existing).Error
}

// GetAll retrieves all salary records ordered by month descending
func (r *SalaryRepository) GetAll() ([]salary.Salary, error) {
	var salaries []salary.Salary
	err := r.db.Order("month DESC").Find(&salaries).Error
	return salaries, err
}

// GetRecent retrieves the most recent salary records (limit specified)
func (r *SalaryRepository) GetRecent(limit int) ([]salary.Salary, error) {
	var salaries []salary.Salary
	err := r.db.Order("month DESC").Limit(limit).Find(&salaries).Error
	return salaries, err
}

// DeleteByMonth deletes salary record for a specific month
func (r *SalaryRepository) DeleteByMonth(month string) error {
	return r.db.Where("month = ?", month).Delete(&salary.Salary{}).Error
}

// GetCurrentMonth retrieves salary for the current month
func (r *SalaryRepository) GetCurrentMonth() (*salary.Salary, error) {
	currentMonth := salary.GetCurrentMonth()
	return r.GetByMonth(currentMonth)
}

// GetMonthsWithSalary retrieves all months that have salary configured
func (r *SalaryRepository) GetMonthsWithSalary() ([]string, error) {
	var months []string
	err := r.db.Model(&salary.Salary{}).
		Select("month").
		Order("month DESC").
		Pluck("month", &months).Error
	return months, err
}

// GetTotalByMonths calculates total salary for multiple months
func (r *SalaryRepository) GetTotalByMonths(months []string) (float64, error) {
	var total float64
	err := r.db.Model(&salary.Salary{}).
		Select("COALESCE(SUM(monthly_amount), 0)").
		Where("month IN ?", months).
		Scan(&total).Error
	return total, err
}
