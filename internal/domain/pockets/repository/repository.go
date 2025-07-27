package repository

import (
	"expenses-api/internal/domain/pockets"
	"expenses-api/internal/infraestructure/client/mysql"
	"fmt"
	"os"
)

const (
	basePathSqlQueries = "sql/pockets"

	fileSqlQueryGetAll     = "GetAll.sql"
	fileSqlQueryGetActives = "GetActives.sql"
	fileSqlQueryGetByID    = "GetByID.sql"
	fileSqlQueryCreate     = "Create.sql"
	fileSqlQueryUpdate     = "Update.sql"
	fileSqlQueryDelete     = "Delete.sql"
)

func GetAll() ([]pockets.Pocket, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetAll))
	if err != nil {
		return nil, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query))
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
			&objPocket.TotalAmount,
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
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetActives))
	if err != nil {
		return nil, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query))
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
			&objPocket.TotalAmount,
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
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetByID))
	if err != nil {
		return pockets.Pocket{}, err
	}

	resultReview, err := mysql.ClientDB.Query(string(query), pocketID)
	if err != nil {
		return pockets.Pocket{}, err
	}

	var pocket pockets.Pocket
	for resultReview.Next() {
		err = resultReview.Scan(
			&pocket.PocketID,
			&pocket.Name,
			&pocket.Status,
			&pocket.TotalAmount,
			&pocket.CreatedAt,
		)
		if err != nil {
			return pockets.Pocket{}, err
		}
	}

	return pocket, nil
}

func Create(newPocket *pockets.Pocket) error {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryCreate))
	if err != nil {
		return err
	}

	newRecord, err := mysql.ClientDB.Exec(string(query), newPocket.Name, newPocket.Status)
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
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryUpdate))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(
		string(query),
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
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryDelete))
	if err != nil {
		return err
	}

	_, err = mysql.ClientDB.Exec(string(query), pocketID)
	if err != nil {
		return err
	}

	return nil
}
