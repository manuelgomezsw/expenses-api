package usecase

import (
	"errors"
	"expenses-api/internal/application/port"
	"expenses-api/internal/domain/pocket"
	"strings"
)

// PocketUseCase handles pocket-related business logic
type PocketUseCase struct {
	pocketRepo port.PocketRepository
}

// NewPocketUseCase creates a new pocket use case instance
func NewPocketUseCase(pocketRepo port.PocketRepository) *PocketUseCase {
	return &PocketUseCase{
		pocketRepo: pocketRepo,
	}
}

// GetAll retrieves all pockets
func (uc *PocketUseCase) GetAll() ([]pocket.Pocket, error) {
	return uc.pocketRepo.GetAll()
}

// GetByID retrieves a pocket by ID
func (uc *PocketUseCase) GetByID(id uint) (*pocket.Pocket, error) {
	if id == 0 {
		return nil, errors.New("pocket ID is required")
	}
	
	return uc.pocketRepo.GetByID(id)
}

// Create creates a new pocket
func (uc *PocketUseCase) Create(name, description string) (*pocket.Pocket, error) {
	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("pocket name is required")
	}
	
	if len(name) > 255 {
		return nil, errors.New("pocket name cannot exceed 255 characters")
	}
	
	// Check if name already exists by trying to get it
	existing, err := uc.pocketRepo.GetByName(name)
	if err == nil && existing != nil {
		return nil, errors.New("pocket with this name already exists")
	}
	
	// Create pocket
	p := &pocket.Pocket{
		Name:        name,
		Description: strings.TrimSpace(description),
	}
	
	if err := uc.pocketRepo.Create(p); err != nil {
		return nil, err
	}
	
	return p, nil
}

// Update updates an existing pocket
func (uc *PocketUseCase) Update(id uint, name, description string) (*pocket.Pocket, error) {
	if id == 0 {
		return nil, errors.New("pocket ID is required")
	}
	
	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("pocket name is required")
	}
	
	if len(name) > 255 {
		return nil, errors.New("pocket name cannot exceed 255 characters")
	}
	
	// Get existing pocket
	existingPocket, err := uc.pocketRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// Check if name already exists (excluding current pocket)
	if existingPocket.Name != name {
		existing, err := uc.pocketRepo.GetByName(name)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("pocket with this name already exists")
		}
	}
	
	// Update pocket
	existingPocket.Name = name
	existingPocket.Description = strings.TrimSpace(description)
	
	if err := uc.pocketRepo.Update(existingPocket); err != nil {
		return nil, err
	}
	
	return existingPocket, nil
}

// Delete deletes a pocket
func (uc *PocketUseCase) Delete(id uint) error {
	if id == 0 {
		return errors.New("pocket ID is required")
	}
	
	// Check if pocket exists
	_, err := uc.pocketRepo.GetByID(id)
	if err != nil {
		return err
	}
	
	// Note: In a real implementation, we should check for associated expenses
	// For now, we'll allow deletion and let the database constraints handle it
	
	return uc.pocketRepo.Delete(id)
}

// GetByName retrieves a pocket by name
func (uc *PocketUseCase) GetByName(name string) (*pocket.Pocket, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("pocket name is required")
	}
	
	return uc.pocketRepo.GetByName(name)
}
