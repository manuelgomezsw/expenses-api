package repository

import (
	"errors"
	"expenses-api/internal/domain/cycles"
	"expenses-api/internal/infraestructure/client/mysql"
	"expenses-api/internal/util/customdate"
	"fmt"
	"os"
)

const (
	basePathSqlQueries = "sql/cycles"

	fileSqlQueryGetAll        = "GetAll.sql"
	fileSqlQueryGetActive     = "GetActive.sql"
	fileSqlQueryGetByID       = "GetByID.sql"
	fileSqlQueryCreate        = "Create.sql"
	fileSqlQueryCreateHistory = "CreateHistory.sql"
	fileSqlQueryUpdate        = "Update.sql"
	fileSqlQueryClose         = "Close.sql"
	fileSqlQueryDelete        = "Delete.sql"
)

func GetAll() ([]cycles.Cycle, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetAll))
	if err != nil {
		return nil, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query))
	if err != nil {
		return nil, err
	}

	var allCycles []cycles.Cycle
	for resultReview.Next() {
		var objCycle cycles.Cycle

		err = resultReview.Scan(
			&objCycle.CycleID,
			&objCycle.PocketName,
			&objCycle.Name,
			&objCycle.Budget,
			&objCycle.DateInit,
			&objCycle.DateEnd,
			&objCycle.Status,
			&objCycle.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		objCycle.DateInit = customdate.SetToNoon(objCycle.DateInit)
		objCycle.DateEnd = customdate.SetToNoon(objCycle.DateEnd)

		allCycles = append(allCycles, objCycle)
	}

	return allCycles, nil
}

func GetActive() ([]cycles.Cycle, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetActive))
	if err != nil {
		return nil, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query))
	if err != nil {
		return nil, err
	}

	var allCycles []cycles.Cycle
	for resultReview.Next() {
		var objCycle cycles.Cycle

		err = resultReview.Scan(
			&objCycle.CycleID,
			&objCycle.PocketName,
			&objCycle.Name,
			&objCycle.Budget,
			&objCycle.DateInit,
			&objCycle.DateEnd,
			&objCycle.Status,
			&objCycle.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		objCycle.DateInit = customdate.SetToNoon(objCycle.DateInit)
		objCycle.DateEnd = customdate.SetToNoon(objCycle.DateEnd)

		allCycles = append(allCycles, objCycle)
	}

	return allCycles, nil
}

func GetByID(cycleID int) (cycles.Cycle, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetByID))
	if err != nil {
		return cycles.Cycle{}, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query), cycleID)
	if err != nil {
		return cycles.Cycle{}, err
	}

	var objCycle cycles.Cycle
	for resultReview.Next() {

		err = resultReview.Scan(
			&objCycle.CycleID,
			&objCycle.PocketID,
			&objCycle.PocketName,
			&objCycle.Name,
			&objCycle.Budget,
			&objCycle.DateInit,
			&objCycle.DateEnd,
			&objCycle.Status,
			&objCycle.CreatedAt,
		)
		if err != nil {
			return cycles.Cycle{}, err
		}
	}

	objCycle.DateInit = customdate.SetToNoon(objCycle.DateInit)
	objCycle.DateEnd = customdate.SetToNoon(objCycle.DateEnd)

	return objCycle, nil
}

func Create(cycle *cycles.Cycle) error {
	// Convertir dateInit a "YYYY-MM-DD HH:MM:SS"
	dateInitFormatted, err := customdate.ParseAndFormatDateMySql(cycle.DateInit)
	if err != nil {
		return errors.New("fecha inicial inválida: " + err.Error())
	}

	// Convertir dateEnd a "YYYY-MM-DD HH:MM:SS"
	dateEndFormatted, err := customdate.ParseAndFormatDateMySql(cycle.DateEnd)
	if err != nil {
		return errors.New("fecha final inválida: " + err.Error())
	}

	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryCreate))
	if err != nil {
		return err
	}

	newRecord, err := mysql.ClientDB.Exec(
		string(query),
		cycle.PocketID,
		cycle.Name,
		cycle.Budget,
		dateInitFormatted,
		dateEndFormatted,
	)
	if err != nil {
		return err
	}

	cycleID, err := newRecord.LastInsertId()
	if err != nil {
		return err
	}
	cycle.CycleID = int(cycleID)

	return nil
}

func CreateHistory(cycleHistory cycles.History) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryCreateHistory))
	if err != nil {
		return err
	}

	// Convertir dateInit a "YYYY-MM-DD HH:MM:SS"
	dateInitFormatted, err := customdate.ParseAndFormatDateMySql(cycleHistory.DateInit)
	if err != nil {
		return errors.New("fecha inicial inválida: " + err.Error())
	}

	// Convertir dateEnd a "YYYY-MM-DD HH:MM:SS"
	dateEndFormatted, err := customdate.ParseAndFormatDateMySql(cycleHistory.DateEnd)
	if err != nil {
		return errors.New("fecha final inválida: " + err.Error())
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		cycleHistory.PocketName,
		cycleHistory.CycleName,
		cycleHistory.Budget,
		cycleHistory.Spent,
		cycleHistory.SpentRatio,
		dateInitFormatted,
		dateEndFormatted,
	)
	if err != nil {
		return err
	}

	return nil
}

func Update(cycleID int, currentCycle *cycles.Cycle) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryUpdate))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		currentCycle.PocketID,
		currentCycle.Budget,
		currentCycle.DateInit,
		currentCycle.DateEnd,
		currentCycle.Status,
		cycleID,
	)
	if err != nil {
		return err
	}

	currentCycle.CycleID = cycleID

	return nil
}

func Close(cycleID int) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryClose))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
		cycleID,
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
