package users_db

import (
	"context"
	"fmt"
	"os"

	"github.com/Sora8d/bookstore_utils-go/logger"

	pgx "github.com/jackc/pgx/v4"
)

const (
	postgres_users_username = "postgres_users_username"
	postgres_users_password = "postgres_users_password"
	postgres_users_host     = "postgres_users_host"
	postgres_users_schema   = "postgres_users_schema"
)

var (
	Client   postGresInterface
	username = os.Getenv(postgres_users_username)
	password = os.Getenv(postgres_users_password)
	host     = os.Getenv(postgres_users_host)
	schema   = os.Getenv(postgres_users_schema)
)

type postGresInterface interface {
	Get
	FindByEmailAndPassword
	FindByStatus
	Delete
	Update
	Save
}

type postGresObject struct {
	conn *pgx.Conn
}

/*
postgres://
postgres://localhost
postgres://localhost:5432
postgres://localhost/mydb
postgres://user@localhost
postgres://user:secret@localhost
postgres://other@localhost/otherdb?connect_timeout=10&application_name=myapp
postgres://localhost/mydb?user=other&password=secret

possible postgresql urls
*/
func init() {
	//When im able to set environment variables on vscode ill update this
	datasourceName := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		username,
		password,
		schema)
	var err error
	newConn, err := pgx.Connect(context.Background(), datasourceName)
	if err != nil {
		logger.Error("Fatal error initializing db", err)
		panic(err)
	}
	if err = newConn.Ping(context.Background()); err != nil {
		logger.Error("Fatal error initializing db", err)
		panic(err)
	}
	Client = &postGresObject{conn: newConn}
	logger.Info("database succesfully configured")
}
