package repository

import (
	"expenses-api/internal/domain/expenses"
	"expenses-api/internal/infraestructure/client/mysql"
	"fmt"
	"os"
)

const (
	basePathSqlQueries = "sql/expenses"

	fileSqlQueryGetByActiveCycles = "GetByActiveCycles.sql"
	fileSqlCreate                 = "Create.sql"
	fileSqlUpdate                 = "Update.sql"
	fileSqlDelete                 = "Delete.sql"
)

func GetByActiveCycles() ([]expenses.Expense, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetByActiveCycles))
	if err != nil {
		return nil, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query))
	if err != nil {
		return nil, err
	}

	var allExpenses []expenses.Expense
	for resultReview.Next() {
		var expense expenses.Expense

		err = resultReview.Scan(
			&expense.ExpenseID,
			&expense.Name,
			&expense.Value,
			&expense.CycleID,
			&expense.CycleName,
			&expense.PocketID,
			&expense.PocketName,
			&expense.PaymentTypeID,
			&expense.PaymentTypeName,
			&expense.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		allExpenses = append(allExpenses, expense)
	}

	return allExpenses, nil
}

func Create(expense *expenses.Expense) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlCreate))
	if err != nil {
		return err
	}

	newRecord, err := mysql.ClientDB.Exec(
		string(query),
		expense.Name,
		expense.Value,
		expense.CycleID,
		expense.PaymentTypeID,
	)
	if err != nil {
		return err
	}

	expenseID, err := newRecord.LastInsertId()
	if err != nil {
		return err
	}
	expense.ExpenseID = int(expenseID)

	return nil
}

func Update(expenseID int, currentExpense *expenses.Expense) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlUpdate))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		currentExpense.Name,
		currentExpense.Value,
		currentExpense.CycleID,
		currentExpense.PaymentTypeID,
		expenseID,
	)
	if err != nil {
		return err
	}

	currentExpense.ExpenseID = expenseID

	return nil
}

func Delete(expenseID int) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlDelete))
	if err != nil {
		return err
	}
	_, err = mysql.ClientDB.Exec(string(query), expenseID)
	if err != nil {
		return err
	}

	return nil
}
