package repository

import (
	"expenses-api/internal/domain/pocket"

	"gorm.io/gorm"
)

// PocketRepository handles pocket-related database operations
type PocketRepository struct {
	*BaseRepository
}

// NewPocketRepository creates a new pocket repository instance
func NewPocketRepository(db *gorm.DB) *PocketRepository {
	return &PocketRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetAll retrieves all pockets ordered by name
func (r *PocketRepository) GetAll() ([]pocket.Pocket, error) {
	var pockets []pocket.Pocket
	err := r.db.Order("name ASC").Find(&pockets).Error
	return pockets, err
}

// GetByID retrieves a pocket by ID
func (r *PocketRepository) GetByID(id uint) (*pocket.Pocket, error) {
	var p pocket.Pocket
	err := r.db.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create creates a new pocket
func (r *PocketRepository) Create(p *pocket.Pocket) error {
	return r.db.Create(p).Error
}

// Update updates an existing pocket
func (r *PocketRepository) Update(p *pocket.Pocket) error {
	return r.db.Save(p).Error
}

// Delete deletes a pocket by ID
func (r *PocketRepository) Delete(id uint) error {
	return r.db.Delete(&pocket.Pocket{}, id).Error
}

// GetByName retrieves a pocket by name
func (r *PocketRepository) GetByName(name string) (*pocket.Pocket, error) {
	var p pocket.Pocket
	err := r.db.Where("name = ?", name).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ExistsByName checks if a pocket with the given name exists
func (r *PocketRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&pocket.Pocket{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcludingID checks if a pocket with the given name exists, excluding a specific ID
func (r *PocketRepository) ExistsByNameExcludingID(name string, id uint) (bool, error) {
	var count int64
	err := r.db.Model(&pocket.Pocket{}).
		Where("name = ? AND id != ?", name, id).
		Count(&count).Error
	return count > 0, err
}

// GetWithFixedExpensesCount retrieves pockets with count of associated fixed expenses
func (r *PocketRepository) GetWithFixedExpensesCount() ([]PocketWithStats, error) {
	var results []PocketWithStats

	err := r.db.Table("pockets p").
		Select(`
			p.id,
			p.name,
			p.description,
			p.created_at,
			COALESCE(fe_stats.fixed_expenses_count, 0) as fixed_expenses_count,
			COALESCE(fe_stats.total_fixed_amount, 0) as total_fixed_amount,
			COALESCE(fe_stats.active_months_count, 0) as active_months_count
		`).
		Joins(`
			LEFT JOIN (
				SELECT 
					pocket_id,
					COUNT(*) as fixed_expenses_count,
					SUM(amount) as total_fixed_amount,
					COUNT(DISTINCT month) as active_months_count
				FROM fixed_expenses
				GROUP BY pocket_id
			) fe_stats ON p.id = fe_stats.pocket_id
		`).
		Order("p.name ASC").
		Scan(&results).Error

	return results, err
}

// CanBeDeleted checks if a pocket can be safely deleted (no associated expenses)
func (r *PocketRepository) CanBeDeleted(id uint) (bool, error) {
	var count int64

	// Check fixed expenses
	err := r.db.Model(&struct {
		TableName string `gorm:"-" sql:"fixed_expenses"`
	}{}).Where("pocket_id = ?", id).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// PocketWithStats represents a pocket with statistics
type PocketWithStats struct {
	ID                 uint    `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	CreatedAt          string  `json:"created_at"`
	FixedExpensesCount int     `json:"fixed_expenses_count"`
	TotalFixedAmount   float64 `json:"total_fixed_amount"`
	ActiveMonthsCount  int     `json:"active_months_count"`
}
