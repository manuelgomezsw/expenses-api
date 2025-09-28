package pocket

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Pocket represents organizational categories for expenses
// Maps to frontend interface: Pocket { id?, name, description?, created_at? }
type Pocket struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for GORM
func (Pocket) TableName() string {
	return "pockets"
}

// BeforeCreate hook to validate and clean data before creation
func (p *Pocket) BeforeCreate(tx *gorm.DB) error {
	return p.validate()
}

// BeforeUpdate hook to validate and clean data before update
func (p *Pocket) BeforeUpdate(tx *gorm.DB) error {
	return p.validate()
}

// validate performs validation and data cleaning
func (p *Pocket) validate() error {
	// Clean and validate name
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return errors.New("pocket name cannot be empty")
	}

	if len(p.Name) > 255 {
		return errors.New("pocket name cannot exceed 255 characters")
	}

	// Clean description
	p.Description = strings.TrimSpace(p.Description)

	return nil
}

// IsEmpty checks if the pocket has no associated expenses
func (p *Pocket) IsEmpty(db *gorm.DB) bool {
	var count int64

	// Check fixed expenses
	db.Table("fixed_expenses").Where("pocket_id = ?", p.ID).Count(&count)
	if count > 0 {
		return false
	}

	return true
}
