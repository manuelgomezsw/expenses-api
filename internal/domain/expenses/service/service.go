package service

import (
	"expenses-api/internal/domain/expenses"
	"expenses-api/internal/domain/expenses/repository"
)

func GetByPocketID(pocketID int) ([]expenses.Expense, error) {
	allExpenses, err := repository.GetByPocketID(pocketID)
	if err != nil {
		return nil, err
	}

	return allExpenses, nil
}

func GetByPaymentTypeID(paymentTypeID int16) ([]expenses.Expense, error) {
	allExpenses, err := repository.GetByPaymentTypeID(paymentTypeID)
	if err != nil {
		return nil, err
	}

	return allExpenses, nil
}

func GetByID(expenseID int) (expenses.Expense, error) {
	expense, err := repository.GetByID(expenseID)
	if err != nil {
		return expenses.Expense{}, err
	}

	return expense, nil
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
