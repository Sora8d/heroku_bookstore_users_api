package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Sora8d/heroku_bookstore_users_api/domain/users"
	"github.com/Sora8d/heroku_bookstore_users_api/services"

	"github.com/Sora8d/bookstore_utils-go/rest_errors"

	"github.com/Sora8d/bookstore_oauth-go/oauth"
	"github.com/gin-gonic/gin"
)

func getUserId(Paramid string) (int64, rest_errors.RestErr) {
	userId, userErr := strconv.ParseInt(Paramid, 10, 64)
	if userErr != nil {
		err := rest_errors.NewBadRequestErr("invalid user id")
		return 0, err
	}
	return userId, nil
}

func Create(c *gin.Context) {
	var user users.User
	fmt.Println(user)
	//First way
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//TODO: Handle error
		return
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		restErr := rest_errors.NewBadRequestErr("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	/*We can also use c.ShouldBindJSON(&user), that replaces everythin from
	Readall() to Unmarshal()*/

	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status(), saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status(), err)
	}

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}
	/* albeit this works, it breaks the way structures are arranged, so the video creates and uses a function from the service package
	reqUser := users.User{Id: userId}
	getErr := reqUser.Get()
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	*/

	reqUser, reqErr := services.UsersService.GetUser(userId)
	if reqErr != nil {
		c.JSON(reqErr.Status(), reqErr)
		return
	}

	/* This is not needed, c.JSON takes care of transforming the struct to JSON, but in your own implementation of a router it will be useful
	retUserJSON, jsonErr := json.Marshal(reqUser)
	if jsonErr != nil {
		//TODO: implement marshal error (Does this need an implementation?)
		return
	}
	*/

	if oauth.GetCallerId(c.Request) == reqUser.Id {
		c.JSON(http.StatusOK, reqUser.Marshall(false))
	}

	c.JSON(http.StatusOK, reqUser.Marshall(oauth.IsPublic(c.Request)))
}

func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestErr("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	user.Id = userId

	IsPartial := c.Request.Method == http.MethodPatch

	result, UpdErr := services.UsersService.UpdateUser(IsPartial, user)
	if UpdErr != nil {
		c.JSON(UpdErr.Status(), UpdErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}
	if err := services.UsersService.DeleteUser(userId); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall((c.GetHeader("X-Public") == "true")))
}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestErr("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	user, err := services.UsersService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, user)
}
