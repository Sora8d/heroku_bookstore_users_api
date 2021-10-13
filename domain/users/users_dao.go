package users

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Sora8d/heroku_bookstore_users_api/datasources/postgresql/users_db"
	"github.com/Sora8d/heroku_bookstore_users_api/utils/date"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/Sora8d/bookstore_utils-go/rest_errors"
)

const (
	erroruniqueconstraint       = "users_email_key"
	queryInsertUser             = "INSERT INTO users (first_name, last_name, email, status, date_created, password) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;"
	queryGetUser                = "SELECT id, first_name, last_name, email, status, date_created FROM users WHERE id = $1;"
	queryUpdateUser             = "UPDATE users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
	queryDeleteUser             = "DELETE FROM users WHERE id=$1;"
	queryFindUserByStatus       = "SELECT id, first_name, last_name, email, status, date_created FROM users WHERE status=$1;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, status, date_Created FROM users WHERE email=$1 AND password=$2;"
)

var usersDB = users_db.Client

func (user *User) Get() rest_errors.RestErr {
	result := usersDB.Get(queryGetUser, user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.DateCreated); err != nil {
		return rest_errors.NewNotFoundError("No matching ids")
	}
	return nil
}

func (user *User) FindByEmailAndPassword() rest_errors.RestErr {
	result, err := usersDB.Query(queryFindByEmailAndPassword, user.Email, user.Password)
	if err != nil {
		//TODO do stuff
		logger.Error("Error in FindByEmailAndPassword function, ", err)
		return rest_errors.NewInternalServerError("There was an unexpected error in the login process", errors.New("database error"))
	}
	defer result.Close()
	if !result.Next() {
		return rest_errors.NewBadRequestErr("Incorrect user credentials")
	}
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.DateCreated); err != nil {
		logger.Error("Error parsing login info in FindByEmailAndPassword function", err)
		resterr := rest_errors.NewInternalServerError("Error trying to validate credentials", errors.New("database error"))
		return resterr
	}

	return nil
}

/*
func (user *User) Get(userId int64) (*User, *errors.RestErr) {
	return nil, nil
}
*/
func (user *User) Save() rest_errors.RestErr {
	user.DateCreated = date.GetNowString()
	row := usersDB.Insert(queryInsertUser, user.FirstName, user.LastName, user.Email, user.Status, user.DateCreated, user.Password)
	var userId int64
	err := row.Scan(&userId)
	if err != nil {
		//This is a placeholder of a postgres_utils that should parse errors accordingly
		if strings.Contains(err.Error(), erroruniqueconstraint) {
			return rest_errors.NewBadRequestErr("email given is invalid")
		}
		resterr := rest_errors.NewInternalServerError("Error validating user", errors.New("database error"))
		return resterr
	}
	user.Id = userId

	return nil
}

func (user *User) Update() rest_errors.RestErr {
	err := usersDB.Execute(queryUpdateUser, user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		resterr := rest_errors.NewInternalServerError("Error updating user info", errors.New("database error"))
		return resterr
	}
	return nil
}

func (user *User) Delete() rest_errors.RestErr {
	err := usersDB.Execute(queryDeleteUser, user.Id)
	if err != nil {
		return rest_errors.NewBadRequestErr("There was an error using method delete with given id")
	}
	return nil
}

func (user *User) FindByStatus(status string) (Users, rest_errors.RestErr) {
	rows, err := usersDB.Query(queryFindUserByStatus, status)
	if err != nil {
		//TODO do stuff
		logger.Error("Error in FindByStatus function", err)
		return nil, rest_errors.NewInternalServerError("There was an error getting the search results", errors.New("database error"))
	}
	defer rows.Close()

	results := make(Users, 0)
	for rows.Next() {
		var current User
		if err := rows.Scan(&current.Id, &current.FirstName, &current.LastName, &current.Email, &current.DateCreated, &current.Status); err != nil {
			logger.Error("Error parsing scans of the function FindByStatus", err)
			return nil, rest_errors.NewInternalServerError("There was an error in one of the results", errors.New("database error"))
		}
		results = append(results, current)
	}
	if len(results) == 0 {
		resterr := rest_errors.NewNotFoundError(fmt.Sprintf("No users matching status %s", status))
		return nil, resterr
	}
	return results, nil
}
