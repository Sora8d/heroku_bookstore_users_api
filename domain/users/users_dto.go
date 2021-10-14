package users

import (
	"strings"

	"github.com/Sora8d/bookstore_utils-go/rest_errors"
)

const (
	StatusActive = "active"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"datecreated"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type Users []User

func (user *User) Validate() rest_errors.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		resterr := rest_errors.NewBadRequestErr("invalid email address")
		return resterr
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		resterr := rest_errors.NewBadRequestErr("invalid password")
		return resterr
	}
	return nil
}
