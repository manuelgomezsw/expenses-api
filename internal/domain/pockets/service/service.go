package service

import (
	"expenses-api/internal/domain/pockets"
	"expenses-api/internal/domain/pockets/repository"
	"fmt"
	"time"
)

func Get() ([]pockets.Pocket, error) {
	pocket, err := repository.Get()
	if err != nil {
		return nil, err
	}

	return pocket, nil
}

func GetByID(pocketID int64) (pockets.Pocket, error) {
	pocket, err := repository.GetByID(pocketID)
	if err != nil {
		return pockets.Pocket{}, err
	}

	return pocket, nil
}

func Create(pocket *pockets.Pocket) error {
	if err := checkDates(pocket.DateInit, pocket.DateEnd); err != nil {
		return err
	}

	if err := repository.Create(pocket); err != nil {
		return err
	}

	return nil
}

func Update(pocketID int64, pocket *pockets.Pocket) error {
	if pocket.DateInit != "" && pocket.DateEnd != "" {
		if err := checkDates(pocket.DateInit, pocket.DateEnd); err != nil {
			return err
		}
	}

	if err := repository.Update(pocketID, pocket); err != nil {
		return err
	}

	return nil
}

func Delete(pocketID int64) error {
	if err := repository.Delete(pocketID); err != nil {
		return err
	}

	return nil
}

func checkDates(dateInit, dateEnd string) error {
	formatDate := "2006-01-02"

	parsedDateInit, err := time.Parse(formatDate, dateInit)
	if err != nil {
		return err
	}

	parsedDateEnd, err := time.Parse(formatDate, dateEnd)
	if err != nil {
		return err
	}

	if parsedDateEnd.After(parsedDateInit) {
		return nil
	}

	return fmt.Errorf("date_end %s must be grather than date_init %s", dateEnd, dateInit)
}
