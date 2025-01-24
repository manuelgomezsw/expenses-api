package repository

import (
	"expenses-api/internal/domain/pockets"
	"expenses-api/internal/infraestructure/client/mysql"
)

func GetAll() ([]pockets.Pocket, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name, status, created_at FROM pockets ORDER BY name")
	if err != nil {
		return nil, err
	}

	var allPockets []pockets.Pocket
	for resultReview.Next() {
		var objPocket pockets.Pocket

		err = resultReview.Scan(
			&objPocket.PocketID,
			&objPocket.Name,
			&objPocket.Status,
			&objPocket.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		allPockets = append(allPockets, objPocket)
	}

	return allPockets, nil
}

func GetActives() ([]pockets.Pocket, error) {
	resultReview, err := mysql.ClientDB.Query(
		"SELECT id, name, status, created_at FROM pockets WHERE status = true ORDER BY name")
	if err != nil {
		return nil, err
	}

	var allPockets []pockets.Pocket
	for resultReview.Next() {
		var objPocket pockets.Pocket

		err = resultReview.Scan(
			&objPocket.PocketID,
			&objPocket.Name,
			&objPocket.Status,
			&objPocket.CreatedAt,
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
		"SELECT id, name, status, created_at FROM pockets WHERE id = ?", pocketID)
	if err != nil {
		return pockets.Pocket{}, err
	}

	var pocket pockets.Pocket
	for resultReview.Next() {
		err = resultReview.Scan(
			&pocket.PocketID,
			&pocket.Name,
			&pocket.Status,
			&pocket.CreatedAt,
		)
		if err != nil {
			return pockets.Pocket{}, err
		}
	}

	return pocket, nil
}

func Create(newPocket *pockets.Pocket) error {
	newRecord, err := mysql.ClientDB.Exec(
		"INSERT INTO pockets (name, status) VALUES (?, ?)",
		newPocket.Name,
		newPocket.Status,
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

func Update(pocketID int64, pocket *pockets.Pocket) error {
	_, err := mysql.ClientDB.Exec(
		"UPDATE pockets SET name = ?, status = ? WHERE id = ?",
		pocket.Name,
		pocket.Status,
		pocket.PocketID,
	)
	if err != nil {
		return err
	}

	pocket.PocketID = pocketID

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
