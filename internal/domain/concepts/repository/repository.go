package repository

import (
	"expenses-api/internal/domain/concepts"
	"expenses-api/internal/infraestructure/client/mysql"
	"fmt"
	"os"
	"time"
)

const (
	basePathSqlQueries = "sql/concepts"

	fileSqlQueryGetByPocketID = "GetByPocketID.sql"
	fileSqlQueryCreate        = "Create.sql"
	fileSqlQueryUpdate        = "Update.sql"
	fileSqlQueryDelete        = "Delete.sql"
)

func GetByPocketID(pocketID int) ([]concepts.Concept, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetByPocketID))
	if err != nil {
		return nil, err
	}

	conceptsByPocketResult, err := mysql.ClientDB.Query(string(query), pocketID)
	if err != nil {
		return nil, err
	}

	var conceptsByPocket []concepts.Concept

	for conceptsByPocketResult.Next() {
		var objConcept concepts.Concept

		err = conceptsByPocketResult.Scan(
			&objConcept.ID,
			&objConcept.Name,
			&objConcept.Value,
			&objConcept.PocketID,
			&objConcept.Payed,
			&objConcept.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		conceptsByPocket = append(conceptsByPocket, objConcept)
	}

	return conceptsByPocket, nil
}

func Create(concept *concepts.Concept) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryCreate))
	if err != nil {
		return err
	}

	newRecord, err := mysql.ClientDB.Exec(
		string(query),
		concept.Name,
		concept.Value,
		concept.PocketID,
		concept.Payed,
		time.Now(),
	)
	if err != nil {
		return err
	}

	conceptID, err := newRecord.LastInsertId()
	if err != nil {
		return err
	}
	concept.ID = int(conceptID)

	return nil
}

func Update(conceptID int, currentConcept *concepts.Concept) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryUpdate))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		currentConcept.Name,
		currentConcept.Value,
		currentConcept.PocketID,
		currentConcept.Payed,
		time.Now(),
		conceptID,
	)
	if err != nil {
		return err
	}

	currentConcept.ID = conceptID

	return nil
}

func Delete(cycleID int) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryDelete))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(string(query), cycleID)
	if err != nil {
		return err
	}

	return nil
}
