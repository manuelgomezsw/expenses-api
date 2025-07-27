package service

import (
	"expenses-api/internal/domain/payments"
	"expenses-api/internal/domain/payments/repository"
)

func Get() ([]payments.Type, error) {
	paymentsType, err := repository.Get()
	if err != nil {
		return nil, err
	}

	return paymentsType, nil
}

func GetByID(paymentTypeID int16) (payments.Type, error) {
	paymentType, err := repository.GetByID(paymentTypeID)
	if err != nil {
		return payments.Type{}, err
	}

	return paymentType, nil
}

func Create(paymentType *payments.Type) error {
	if err := repository.Create(paymentType); err != nil {
		return err
	}

	return nil
}

func Update(paymentTypeID int16, paymentType *payments.Type) error {
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
