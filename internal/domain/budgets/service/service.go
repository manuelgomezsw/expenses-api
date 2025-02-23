package service

import (
	"expenses-api/internal/domain/budgets"
	cycleRepository "expenses-api/internal/domain/cycles/repository"
	"expenses-api/internal/domain/expenses"
	expenseRepository "expenses-api/internal/domain/expenses/repository"
)

func Calculate(cycleID int) (budgets.Budget, error) {
	budget := budgets.Budget{}

	currentCycle, err := cycleRepository.GetByID(cycleID)
	if err != nil {
		return budgets.Budget{}, err
	}
	budget.CycleID = cycleID
	budget.Budget = currentCycle.Budget

	currentExpenses, err := expenseRepository.GetByCycleID(cycleID)
	if err != nil {
		return budgets.Budget{}, err
	}

	budget.Spent = SumExpenses(currentExpenses)
	budget.SpentRatio = GetSpentRatio(float64(currentCycle.Budget), float64(budget.Spent))

	return budget, nil
}

func SumExpenses(currentExpenses []expenses.Expense) int64 {
	spent := int64(0)
	for _, expense := range currentExpenses {
		spent += expense.Value
	}

	return spent
}

func GetSpentRatio(budget, spent float64) int {
	return int((spent / budget) * 100)
}
