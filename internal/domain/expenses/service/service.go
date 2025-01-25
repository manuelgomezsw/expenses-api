package service

import (
	"expenses-api/internal/domain/expenses"
	"expenses-api/internal/domain/expenses/repository"
)

func GetByActiveCycles() ([]expenses.Expense, error) {
	allExpenses, err := repository.GetByActiveCycles()
	if err != nil {
		return nil, err
	}

	return allExpenses, nil
}

func Create(expense *expenses.Expense) error {
	if err := repository.Create(expense); err != nil {
		return err
	}

	return nil
}

func Update(expenseID int, expense *expenses.Expense) error {
	if err := repository.Update(expenseID, expense); err != nil {
		return err
	}

	return nil
}

func Delete(expenseID int) error {
	if err := repository.Delete(expenseID); err != nil {
		return err
	}

	return nil
}
