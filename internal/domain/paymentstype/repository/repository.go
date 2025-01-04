package repository

import (
	"expenses-api/internal/domain/paymentstype"
	"expenses-api/internal/infraestructure/client/mysql"
)

func Get() ([]paymentstype.PaymentType, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name FROM payment_types ORDER BY id")
	if err != nil {
		return nil, err
	}

	var paymentsType []paymentstype.PaymentType
	for resultReview.Next() {
		var paymentType paymentstype.PaymentType

		err = resultReview.Scan(&paymentType.PaymentTypeID, &paymentType.Name)
		if err != nil {
			return nil, err
		}

		paymentsType = append(paymentsType, paymentType)
	}

	return paymentsType, nil
}

func GetByID(paymentTypeID int16) (paymentstype.PaymentType, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name FROM payment_types WHERE id = ?", paymentTypeID)
	if err != nil {
		return paymentstype.PaymentType{}, err
	}

	var paymentType paymentstype.PaymentType
	for resultReview.Next() {
		err = resultReview.Scan(&paymentType.PaymentTypeID, &paymentType.Name)
		if err != nil {
			return paymentstype.PaymentType{}, err
		}
	}

	return paymentType, nil
}

func Create(newPaymentType *paymentstype.PaymentType) error {
	_, err := mysql.ClientDB.Exec(
		"INSERT INTO payment_types (id, name) VALUES (?, ?)",
		newPaymentType.PaymentTypeID,
		newPaymentType.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func Update(paymentTypeID int16, newPaymentType *paymentstype.PaymentType) error {
	_, err := mysql.ClientDB.Exec(
		"UPDATE payment_types SET name = ? WHERE id = ?",
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
	_, err := mysql.ClientDB.Exec(
		"DELETE FROM payment_types WHERE id = ?",
		paymentTypeID,
	)
	if err != nil {
		return err
	}

	return nil
}
