package repository

import (
	"expenses-api/internal/domain/expenses"
	"expenses-api/internal/infraestructure/client/mysql"
	"fmt"
	"os"
	"strings"
)

const (
	basePathSqlQueries = "sql/expenses"

	fileSqlQueryGetByActiveCycles = "GetByActiveCycles.sql"
	fileSqlCreate                 = "Create.sql"
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

func GetByID(expenseID int) (expenses.Expense, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name, value, pocket_id, payment_type_id FROM expenses WHERE id = ?", expenseID)
	if err != nil {
		return expenses.Expense{}, err
	}

	var expense expenses.Expense
	for resultReview.Next() {
		err = resultReview.Scan(
			&expense.ExpenseID,
			&expense.Name,
			&expense.Value,
			&expense.PocketID,
			&expense.PaymentTypeID,
		)
		if err != nil {
			return expenses.Expense{}, err
		}
	}

	return expense, nil
}

func GetByPocketID(pocketID int) ([]expenses.Expense, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name, value, pocket_id, payment_type_id FROM expenses WHERE pocket_id = ?",
		pocketID,
	)
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
			&expense.PocketID,
			&expense.PaymentTypeID,
		)
		if err != nil {
			return nil, err
		}

		allExpenses = append(allExpenses, expense)
	}

	return allExpenses, nil
}

func GetByPaymentTypeID(paymentTypeID int16) ([]expenses.Expense, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name, value, pocket_id, payment_type_id FROM expenses WHERE payment_type_id = ?",
		paymentTypeID,
	)
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
			&expense.PocketID,
			&expense.PaymentTypeID,
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
	query, params := buildQueryUpdate(expenseID, currentExpense)
	_, err := mysql.ClientDB.Exec(query, params...)
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

func buildQueryUpdate(expenseID int, newPocket *expenses.Expense) (string, []interface{}) {
	query := "UPDATE expenses SET "
	var params []interface{}

	if newPocket.Name != "" {
		query += "name = ?, "
		params = append(params, newPocket.Name)
	}
	if newPocket.Value != 0 {
		query += "value = ?, "
		params = append(params, newPocket.Value)
	}
	if newPocket.CycleID != 0 {
		query += "cycle_id = ?, "
		params = append(params, newPocket.PocketID)
	}
	if newPocket.PaymentTypeID != 0 {
		query += "payment_type_id = ?, "
		params = append(params, newPocket.PaymentTypeID)
	}

	if len(params) == 0 {
		return "", nil
	}

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = ?"
	params = append(params, expenseID)

	return query, params
}
