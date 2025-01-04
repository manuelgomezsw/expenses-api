package service

import (
	"expenses-api/internal/domain/paymentstype"
	"expenses-api/internal/domain/paymentstype/repository"
)

func Get() ([]paymentstype.PaymentType, error) {
	paymentsType, err := repository.Get()
	if err != nil {
		return nil, err
	}

	return paymentsType, nil
}

func GetByID(paymentTypeID int16) (paymentstype.PaymentType, error) {
	paymentType, err := repository.GetByID(paymentTypeID)
	if err != nil {
		return paymentstype.PaymentType{}, err
	}

	return paymentType, nil
}

func Create(paymentType *paymentstype.PaymentType) error {
	if err := repository.Create(paymentType); err != nil {
		return err
	}

	return nil
}

func Update(paymentTypeID int16, paymentType *paymentstype.PaymentType) error {
	if err := repository.Update(paymentTypeID, paymentType); err != nil {
		return err
	}

	return nil
}

func Delete(paymentTypeID int16) error {
	if err := repository.Delete(paymentTypeID); err != nil {
		return err
	}

	return nil
}
