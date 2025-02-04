package service

import (
	"expenses-api/internal/domain/concepts"
	"expenses-api/internal/domain/concepts/repository"
)

func GetByPocketID(pocketID int) ([]concepts.Concept, error) {
	return repository.GetByPocketID(pocketID)
}

func Create(concept *concepts.Concept) error {
	return repository.Create(concept)
}

func Update(conceptID int, currentConcept *concepts.Concept) error {
	return repository.Update(conceptID, currentConcept)
}

func Delete(cycleID int) error {
	return repository.Delete(cycleID)
}
