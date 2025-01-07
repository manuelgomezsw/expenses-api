package repository

import (
	"expenses-api/internal/domain/pockets"
	"expenses-api/internal/infraestructure/client/mysql"
	"strings"
)

func Get() ([]pockets.Pocket, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name, month, budget, date_init, date_end FROM pockets")
	if err != nil {
		return nil, err
	}

	var allPockets []pockets.Pocket
	for resultReview.Next() {
		var objPocket pockets.Pocket

		err = resultReview.Scan(
			&objPocket.PocketID,
			&objPocket.Name,
			&objPocket.Month,
			&objPocket.Budget,
			&objPocket.DateInit,
			&objPocket.DateEnd,
		)
		if err != nil {
			return nil, err
		}

		allPockets = append(allPockets, objPocket)
	}

	return allPockets, nil
}

func GetByID(pocketID int64) (pockets.Pocket, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name, month, budget, date_init, date_end FROM pockets WHERE id = ?", pocketID)
	if err != nil {
		return pockets.Pocket{}, err
	}

	var pocket pockets.Pocket
	for resultReview.Next() {
		err = resultReview.Scan(
			&pocket.PocketID,
			&pocket.Name,
			&pocket.Month,
			&pocket.Budget,
			&pocket.DateInit,
			&pocket.DateEnd,
		)
		if err != nil {
			return pockets.Pocket{}, err
		}
	}

	return pocket, nil
}

func Create(newPocket *pockets.Pocket) error {
	newRecord, err := mysql.ClientDB.Exec(
		"INSERT INTO pockets (name, month, budget, date_init, date_end) VALUES (?, ?, ?, ?, ?)",
		newPocket.Name,
		newPocket.Month,
		newPocket.Budget,
		newPocket.DateInit,
		newPocket.DateEnd,
	)
	if err != nil {
		return err
	}

	newPocket.PocketID, err = newRecord.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func Update(pocketID int64, newPocket *pockets.Pocket) error {
	query, params := buildQueryUpdate(pocketID, newPocket)
	_, err := mysql.ClientDB.Exec(query, params...)
	if err != nil {
		return err
	}

	newPocket.PocketID = pocketID

	return nil
}

func Delete(pocketID int64) error {
	_, err := mysql.ClientDB.Exec(
		"DELETE FROM pockets WHERE id = ?",
		pocketID,
	)
	if err != nil {
		return err
	}

	return nil
}

func buildQueryUpdate(pocketID int64, newPocket *pockets.Pocket) (string, []interface{}) {
	query := "UPDATE pockets SET "
	params := []interface{}{}

	if newPocket.Name != "" {
		query += "name = ?, "
		params = append(params, newPocket.Name)
	}
	if newPocket.Month != "" {
		query += "month = ?, "
		params = append(params, newPocket.Month)
	}
	if newPocket.Budget != 0 {
		query += "budget = ?, "
		params = append(params, newPocket.Budget)
	}
	if newPocket.DateInit != "" {
		query += "date_init = ?, "
		params = append(params, newPocket.DateInit)
	}
	if newPocket.DateEnd != "" {
		query += "date_end = ?, "
		params = append(params, newPocket.DateEnd)
	}

	if len(params) == 0 {
		return "", nil
	}

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = ?"
	params = append(params, pocketID)

	return query, params
}
