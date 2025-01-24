package repository

import (
	"errors"
	"expenses-api/internal/domain/cycles"
	"expenses-api/internal/infraestructure/client/mysql"
	"expenses-api/internal/util/customdate"
	"fmt"
	"os"
	"strings"
)

const (
	basePathSqlQueries = "sql/cycles"

	fileSqlQueryGetAll    = "GetAll.sql"
	fileSqlQueryGetActive = "GetActive.sql"
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
	resultReview, err := mysql.ClientDB.Query(
		"SELECT "+
			"c.id,"+
			"c.pocket_id,"+
			"p.name,"+
			"c.name,"+
			"c.budget,"+
			"c.date_init,"+
			"c.date_end,"+
			"c.status,"+
			"c.created_at "+
			"FROM cycles c "+
			"JOIN pockets p ON c.pocket_id = p.id "+
			"WHERE c.id = ?",
		cycleID,
	)
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

	newRecord, err := mysql.ClientDB.Exec(
		"INSERT INTO cycles (pocket_id, name, budget, date_init, date_end) VALUES (?, ?, ?, ?, ?)",
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

func Update(cycleID int, currentCycle *cycles.Cycle) error {
	query, params := buildQueryUpdate(cycleID, currentCycle)
	_, err := mysql.ClientDB.Exec(query, params...)
	if err != nil {
		return err
	}

	currentCycle.CycleID = cycleID

	return nil
}

func Delete(cycleID int) error {
	_, err := mysql.ClientDB.Exec("DELETE FROM cycles WHERE id = ?", cycleID)
	if err != nil {
		return err
	}

	return nil
}

func buildQueryUpdate(cycleID int, currentCycle *cycles.Cycle) (string, []interface{}) {
	query := "UPDATE cycles SET "
	var params []interface{}

	if currentCycle.PocketID > 0 {
		query += "pocket_id = ?, "
		params = append(params, currentCycle.PocketID)
	}
	if currentCycle.Budget > 0 {
		query += "budget = ?, "
		params = append(params, currentCycle.Budget)
	}
	if currentCycle.DateInit != "" {
		query += "date_init = ?, "
		params = append(params, currentCycle.DateInit)
	}
	if currentCycle.DateEnd != "" {
		query += "date_end = ?, "
		params = append(params, currentCycle.DateEnd)
	}

	query += "status = ?, "
	params = append(params, currentCycle.Status)

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = ?"
	params = append(params, cycleID)

	return query, params
}
