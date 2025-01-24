package service

import (
	"expenses-api/internal/domain/pockets"
	"expenses-api/internal/domain/pockets/repository"
)

func GetAll() ([]pockets.Pocket, error) {
	pocket, err := repository.GetAll()
	if err != nil {
		return nil, err
	}

	return pocket, nil
}

func GetActives() ([]pockets.Pocket, error) {
	pocket, err := repository.GetActives()
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
	if err := repository.Create(pocket); err != nil {
		return err
	}

	return nil
}

func Update(pocketID int64, pocket *pockets.Pocket) error {
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
