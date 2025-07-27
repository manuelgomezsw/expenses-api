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

	fileSqlQueryGetByID         = "GetByID.sql"
	fileSqlQueryGetByPocketID   = "GetByPocketID.sql"
	fileSqlQueryCreate          = "Create.sql"
	fileSqlQueryUpdate          = "Update.sql"
	fileSqlQueryPayedUpdate     = "PayedUpdate.sql"
	fileSqlQueryDelete          = "Delete.sql"
	fileSqlQueryBulkUpdatePayed = "BulkUpdatePayed.sql"
)

func GetByID(pocketID int) ([]concepts.Concept, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetByID))
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
			&objConcept.PocketName,
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
			&objConcept.PocketName,
			&objConcept.Payed,
			&objConcept.UpdatedAt,
			&objConcept.PaymentDay,
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
		concept.PaymentDay,
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
		currentConcept.PaymentDay,
		conceptID,
	)
	if err != nil {
		return err
	}

	currentConcept.ID = conceptID

	return nil
}

func PayedUpdate(conceptID int, payed *bool) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryPayedUpdate))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		payed,
		time.Now(),
		conceptID,
	)
	if err != nil {
		return err
	}

	return nil
}

func BulkUpdatePayed(pocketID int, payed bool) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryBulkUpdatePayed))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		payed,
		time.Now().Format("2006-01-02"),
		pocketID,
	)
	if err != nil {
		return err
	}

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
