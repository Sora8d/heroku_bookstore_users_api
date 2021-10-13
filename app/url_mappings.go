package app

import (
	"github.com/Sora8d/heroku_bookstore_users_api/controllers/ping"
	"github.com/Sora8d/heroku_bookstore_users_api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/users", users.Create)
	router.GET("/internal/users/search", users.Search)
	router.GET("/users/:user_id", users.Get)
	router.PUT("/users/:user_id", users.Update)
	router.PATCH("/users/:user_id", users.Update)
	router.DELETE("/users/:user_id", users.Delete)
	router.POST("/users/login", users.Login)
}
