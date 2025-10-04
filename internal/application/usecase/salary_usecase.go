package usecase

import (
	"errors"
	"expenses-api/internal/application/port"
	"expenses-api/internal/domain/salary"
	"time"
)

// SalaryUseCase handles salary-related business logic
type SalaryUseCase struct {
	salaryRepo port.SalaryRepository
}

// NewSalaryUseCase creates a new salary use case instance
func NewSalaryUseCase(salaryRepo port.SalaryRepository) *SalaryUseCase {
	return &SalaryUseCase{
		salaryRepo: salaryRepo,
	}
}

// GetByMonth retrieves salary configuration for a specific month
func (uc *SalaryUseCase) GetByMonth(month string) (*salary.Salary, error) {
	if month == "" {
		return nil, errors.New("month is required")
	}

	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return nil, errors.New("invalid month format, must be YYYY-MM")
	}

	return uc.salaryRepo.GetByMonth(month)
}

// GetByMonthWithInheritance obtiene el salario de un mes, heredando del anterior si no existe
func (uc *SalaryUseCase) GetByMonthWithInheritance(month string) (*salary.Salary, error) {
	if month == "" {
		return nil, errors.New("month is required")
	}

	// Validate month format
	date, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, errors.New("invalid month format, must be YYYY-MM")
	}

	// Intentar obtener configuración del mes actual
	currentSalary, err := uc.salaryRepo.GetByMonth(month)
	if err == nil {
		return currentSalary, nil
	}

	// Si no existe, buscar mes anterior
	previousMonth := date.AddDate(0, -1, 0).Format("2006-01")
	previousSalary, err := uc.salaryRepo.GetByMonth(previousMonth)
	if err != nil {
		// No hay configuración anterior, retornar error para que handler use valores por defecto
		return nil, errors.New("no configuration found")
	}

	// Heredar configuración adaptando el mes
	inheritedSalary := &salary.Salary{
		MonthlyAmount: previousSalary.MonthlyAmount,
		Month:         month, // Actualizar al mes solicitado
	}

	return inheritedSalary, nil
}

// GetCurrentMonth retrieves salary for the current month
func (uc *SalaryUseCase) GetCurrentMonth() (*salary.Salary, error) {
	currentMonth := salary.GetCurrentMonth()
	return uc.salaryRepo.GetByMonth(currentMonth)
}

// UpdateSalary updates or creates salary configuration for a month
func (uc *SalaryUseCase) UpdateSalary(monthlyAmount float64, month string) error {
	if month == "" {
		return errors.New("month is required")
	}

	// Validate month format
	if _, err := time.Parse("2006-01", month); err != nil {
		return errors.New("invalid month format, must be YYYY-MM")
	}

	if monthlyAmount < 0 {
		return errors.New("monthly amount cannot be negative")
	}

	salaryConfig := &salary.Salary{
		MonthlyAmount: monthlyAmount,
		Month:         month,
	}

	return uc.salaryRepo.CreateOrUpdate(salaryConfig)
}
