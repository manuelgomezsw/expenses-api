package repository

import (
	"expenses-api/internal/domain/payments"
	"expenses-api/internal/infraestructure/client/mysql"
	"fmt"
	"os"
)

const (
	basePathSqlQueries = "sql/paymentstype"

	fileSqlQueryGet     = "GetAll.sql"
	fileSqlQueryGetByID = "GetByID.sql"
	fileSqlQueryCreate  = "Create.sql"
	fileSqlQueryUpdate  = "Update.sql"
	fileSqlQueryDelete  = "Delete.sql"
)

func Get() ([]payments.Type, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGet))
	if err != nil {
		return nil, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query))
	if err != nil {
		return nil, err
	}

	var paymentsType []payments.Type
	for resultReview.Next() {
		var paymentType payments.Type

		err = resultReview.Scan(&paymentType.PaymentTypeID, &paymentType.Name)
		if err != nil {
			return nil, err
		}

		paymentsType = append(paymentsType, paymentType)
	}

	return paymentsType, nil
}

func GetByID(paymentTypeID int16) (payments.Type, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetByID))
	if err != nil {
		return payments.Type{}, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query), paymentTypeID)
	if err != nil {
		return payments.Type{}, err
	}

	var paymentType payments.Type
	for resultReview.Next() {
		err = resultReview.Scan(&paymentType.PaymentTypeID, &paymentType.Name)
		if err != nil {
			return payments.Type{}, err
		}
	}

	return paymentType, nil
}

func Create(newPaymentType *payments.Type) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryCreate))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		newPaymentType.PaymentTypeID,
		newPaymentType.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func Update(paymentTypeID int16, newPaymentType *payments.Type) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryUpdate))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		newPaymentType.Name,
		paymentTypeID,
	)
	if err != nil {
		return err
	}

	newPaymentType.PaymentTypeID = paymentTypeID

	return nil
}

func Delete(paymentTypeID int16) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryDelete))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(string(query), paymentTypeID)
	if err != nil {
		return err
	}

	return nil
}
