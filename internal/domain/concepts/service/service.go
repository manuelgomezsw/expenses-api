package service

import (
	"expenses-api/internal/domain/concepts"
	"expenses-api/internal/domain/concepts/repository"
)

func GetByID(conceptID int) ([]concepts.Concept, error) {
	return repository.GetByID(conceptID)
}

func GetByPocketID(pocketID int) ([]concepts.Concept, error) {
	return repository.GetByPocketID(pocketID)
}

func Create(concept *concepts.Concept) error {
	return repository.Create(concept)
}

func Update(conceptID int, currentConcept *concepts.Concept) error {
	return repository.Update(conceptID, currentConcept)
}

func PayedUpdate(conceptID int, payed *bool) error {
	return repository.PayedUpdate(conceptID, payed)
}

func Delete(cycleID int) error {
	return repository.Delete(cycleID)
}
