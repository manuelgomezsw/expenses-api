package repository

import (
	"gorm.io/gorm"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB returns the underlying GORM database instance
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}

// Transaction executes a function within a database transaction
func (r *BaseRepository) Transaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// Create creates a new record
func (r *BaseRepository) Create(entity interface{}) error {
	return r.db.Create(entity).Error
}

// Update updates an existing record
func (r *BaseRepository) Update(entity interface{}) error {
	return r.db.Save(entity).Error
}

// Delete deletes a record by ID
func (r *BaseRepository) Delete(entity interface{}, id interface{}) error {
	return r.db.Delete(entity, id).Error
}

// FindByID finds a record by ID
func (r *BaseRepository) FindByID(entity interface{}, id interface{}) error {
	return r.db.First(entity, id).Error
}

// FindAll finds all records
func (r *BaseRepository) FindAll(entities interface{}) error {
	return r.db.Find(entities).Error
}

// Count counts records matching the given conditions
func (r *BaseRepository) Count(model interface{}, count *int64, conditions ...interface{}) error {
	query := r.db.Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Count(count).Error
}

// Exists checks if a record exists with the given conditions
func (r *BaseRepository) Exists(model interface{}, conditions ...interface{}) (bool, error) {
	var count int64
	err := r.Count(model, &count, conditions...)
	return count > 0, err
}

// Paginate applies pagination to a query
func (r *BaseRepository) Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		if pageSize <= 0 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
