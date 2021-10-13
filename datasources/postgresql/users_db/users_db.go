package users_db

import (
	"context"
	"fmt"
	"os"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/joho/godotenv"

	pgx "github.com/jackc/pgx/v4"
)

//env vars
/*
const (
	postgres_users_username = "postgres_users_username"
	postgres_users_password = "postgres_users_password"
	postgres_users_host     = "postgres_users_host"
	postgres_users_schema   = "postgres_users_schema"
)*/

//statements

var (
	Client postGresInterface
)

type postGresInterface interface {
	Get(string, int64) pgx.Row
	Query(string, ...interface{}) (pgx.Rows, error)
	Execute(string, ...interface{}) error
	Insert(string, ...interface{}) pgx.Row
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
	err := godotenv.Load("db_envs.env")
	if err != nil {
		logger.Error("Error loading environment variables", err)
		panic(err)
	}
	//When im able to set environment variables on vscode ill update this
	datasourceName := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		os.Getenv("postgres_users_username"),
		os.Getenv("postgres_users_password"),
		os.Getenv("postgres_users_schema"))
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

type postGresObject struct {
	conn *pgx.Conn
}

func (pgc postGresObject) Get(query string, id int64) pgx.Row {
	row := pgc.conn.QueryRow(context.Background(), query, id)
	return row
}

func (pgc postGresObject) Query(query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := pgc.conn.Query(context.Background(), query, args...)
	return rows, err
}
func (pgc postGresObject) Execute(query string, args ...interface{}) error {
	_, err := pgc.conn.Exec(context.Background(), query, args...)
	return err
}

func (pgc postGresObject) Insert(query string, args ...interface{}) pgx.Row {
	row := pgc.conn.QueryRow(context.Background(), query, args...)
	return row
}
