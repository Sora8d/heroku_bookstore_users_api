package mysql_utils

import (
	"errors"
	"strings"

	"github.com/Sora8d/bookstore_utils-go/rest_errors"

	"github.com/go-sql-driver/mysql"
)

const (
	errorNoRows       = "no rows in result set"
	errorMailRepeated = "users.email"
)

func ParseError(err error) rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			resterr := rest_errors.NewNotFoundError("no record matching given id")
			return resterr
		}
		resterr := rest_errors.NewInternalServerError("error parsing database response", errors.New("database error"))
		return resterr
	}
	switch sqlErr.Number {
	case 1062:
		if strings.Contains(sqlErr.Message, errorMailRepeated) {
			resterr := rest_errors.NewBadRequestErr("email already registered")
			return resterr
		}
		resterr := rest_errors.NewBadRequestErr("invalid data")
		return resterr
	}
	resterr := rest_errors.NewInternalServerError("error processing request", errors.New("database error"))
	return resterr
}
