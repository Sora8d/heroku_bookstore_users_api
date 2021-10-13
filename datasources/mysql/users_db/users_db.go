package users_db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Sora8d/bookstore_utils-go/logger"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysql_users_username = "mysql_users_username"
	mysql_users_password = "mysql_users_password"
	mysql_users_host     = "mysql_users_host"
	mysql_users_schema   = "mysql_users_schema"
)

var (
	Client   *sql.DB
	username = os.Getenv(mysql_users_username)
	password = os.Getenv(mysql_users_password)
	host     = os.Getenv(mysql_users_host)
	schema   = os.Getenv(mysql_users_schema)
)

func init() {
	//When im able to set environment variables on vscode ill update this
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		"test",
		"123",
		"127.0.0.1",
		"users_db")
	var err error
	Client, err = sql.Open("mysql", datasourceName)
	if err != nil {
		logger.Error("Fatal error initializing db", err)
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		logger.Error("Fatal error initializing db", err)
		panic(err)
	}
	logger.Info("database succesfully configured")
}
