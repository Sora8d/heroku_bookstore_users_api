package services

import (
	"github.com/Sora8d/heroku_bookstore_users_api/domain/users"
	"github.com/Sora8d/heroku_bookstore_users_api/utils/crypto_utils"

	errors "github.com/Sora8d/bookstore_utils-go/rest_errors"
)

//This is the solution to services function (lesson 18)
var UsersService usersServiceInterface = &usersService{}

type usersService struct{}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, errors.RestErr)
	CreateUser(users.User) (*users.User, errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, errors.RestErr)
	DeleteUser(int64) errors.RestErr
	SearchUser(string) (users.Users, errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, errors.RestErr)
}

//Here ends the solution
func (s *usersService) GetUser(userId int64) (*users.User, errors.RestErr) {
	reqUser := users.User{Id: userId}
	if err := reqUser.Get(); err != nil {
		return nil, err
	}
	return &reqUser, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) UpdateUser(ispartial bool, user users.User) (*users.User, errors.RestErr) {
	current, err := s.GetUser(user.Id)
	if err != nil {
		return nil, err
	}
	//Now the instructor ties to use User.Validate(), but this break partial
	//When altering the db, for some reason not nul doesnt do anything
	if ispartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}
	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) DeleteUser(userId int64) errors.RestErr {
	user := users.User{Id: userId}
	return user.Delete()
}

func (s *usersService) SearchUser(status string) (users.Users, errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *usersService) LoginUser(request users.LoginRequest) (*users.User, errors.RestErr) {
	dao := &users.User{
		Email:    request.Email,
		Password: crypto_utils.GetMd5(request.Password),
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
