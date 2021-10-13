

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/Sora8d/bookstore_utils-go/rest_errors"
	"github.com/Sora8d/heroku_bookstore_users_api/datasources/mysql/users_db"
	"github.com/Sora8d/heroku_bookstore_users_api/utils/date"
	"github.com/Sora8d/heroku_bookstore_users_api/utils/mysql_utils"
)

const (
	errorNoRows                 = "no rows in result set"
	queryInsertUser             = "INSERT INTO users (first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, status, date_created FROM users WHERE id = ?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus       = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_Created, status FROM users WHERE email=? AND password=?;"
)

var usersDB = users_db.Client

func (user *User) Get() rest_errors.RestErr {
	if err := usersDB.Ping(); err != nil {
		panic(err)
	}
	stmt, err := usersDB.Prepare(queryGetUser)
	if err != nil {
		//This logger we should do with everywhere we have an error
		logger.Error("error when trying to prepare get user statement", err)
		resterr := rest_errors.NewInternalServerError("error trying to get user", errors.New("database error"))
		return resterr
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.DateCreated); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByEmailAndPassword() rest_errors.RestErr {
	if err := usersDB.Ping(); err != nil {
		panic(err)
	}
	stmt, err := usersDB.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		//This logger we should do with everywhere we have an error
		logger.Error("error when trying to prepare get user statement", err)
		resterr := rest_errors.NewInternalServerError("Error trying to validate credentials", errors.New("database error"))
		return resterr
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			resterr := rest_errors.NewBadRequestErr("Incorrect user credentials")
			return resterr
		}
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
	stmt, err := usersDB.Prepare(queryInsertUser)
	if err != nil {
		resterr := rest_errors.NewInternalServerError("Error validating user", errors.New("database error"))
		return resterr
	}
	defer stmt.Close()
	user.DateCreated = date.GetNowString()
	inserResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userId, err := inserResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	user.Id = userId

	return nil
}

func (user *User) Update() rest_errors.RestErr {
	stmt, err := usersDB.Prepare(queryUpdateUser)
	if err != nil {
		resterr := rest_errors.NewInternalServerError("Error updating user info", errors.New("database error"))
		return resterr
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() rest_errors.RestErr {
	stmt, err := usersDB.Prepare(queryDeleteUser)
	if err != nil {
		resterr := rest_errors.NewInternalServerError("Error updating user info", errors.New("database error"))
		return resterr
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) (Users, rest_errors.RestErr) {
	stmt, err := usersDB.Prepare(queryFindUserByStatus)
	if err != nil {
		resterr := rest_errors.NewInternalServerError("Error finding user", errors.New("database error"))
		return nil, resterr
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		resterr := rest_errors.NewInternalServerError("Error finding user", errors.New("database error"))
		return nil, resterr
	}
	defer rows.Close()
	results := make(Users, 0)
	for rows.Next() {
		var current User
		if err := rows.Scan(&current.Id, &current.FirstName, &current.LastName, &current.Email, &current.DateCreated, &current.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, current)
	}
	if len(results) == 0 {
		resterr := rest_errors.NewNotFoundError(fmt.Sprintf("No users matching status %s", status))
		return nil, resterr
	}
	return results, nil
}
