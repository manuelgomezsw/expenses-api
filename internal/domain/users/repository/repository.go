package repository

import (
	"expenses-api/internal/domain/users"
	"expenses-api/internal/infraestructure/client/mysql"
	"expenses-api/internal/util/customdate"
	"fmt"
	"github.com/jinzhu/copier"
	"log"
	"os"
)

const (
	basePathSqlQueries = "sql/users"

	fileSqlQueryCreate      = "Create.sql"
	fileSqlQueryUpdate      = "Update.sql"
	fileSqlQueryDelete      = "Delete.sql"
	fileSqlQueryGetPassword = "GetPassword.sql"
)

func GetByID(userID int) (users.UserOutput, error) {
	return users.UserOutput{}, nil
}

func GetPasswordByUsername(username string) (string, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryGetPassword))
	if err != nil {
		return "", err
	}

	result, err := mysql.ClientDB.Query(string(query), username)
	if err != nil {
		return "", err
	}

	var storedPassword string
	for result.Next() {
		err = result.Scan(&storedPassword)
		if err != nil {
			return storedPassword, err
		}
	}

	return storedPassword, nil
}

func GetByUsername(username string) (users.UserOutput, error) {
	return users.UserOutput{}, nil
}

func Create(user users.UserInput) (users.UserOutput, error) {
	query, err := os.ReadFile(fmt.Sprintf("%s/%s", basePathSqlQueries, fileSqlQueryCreate))
	if err != nil {
		return users.UserOutput{}, err
	}

	newRecord, err := mysql.ClientDB.Exec(
		string(query),
		user.Username,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Email,
		user.CreatedAt,
	)
	if err != nil {
		return users.UserOutput{}, err
	}

	userID, err := newRecord.LastInsertId()
	if err != nil {
		return users.UserOutput{}, err
	}
	user.ID = int(userID)

	var userOutput users.UserOutput
	err = copier.Copy(&userOutput, &user)
	if err != nil {
		log.Fatalf("Copying users entities error: %v", err)
	}
	userOutput.CreatedAt = customdate.SetToNoon(user.CreatedAt)

	return userOutput, nil
}

func Update(userID int, user users.UserInput) (users.UserOutput, error) {
	return users.UserOutput{}, nil
}
