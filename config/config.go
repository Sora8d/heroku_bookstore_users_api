package config

import (
	"os"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/joho/godotenv"
)

type config map[string]string

var Config config

func init() {
	if err := godotenv.Load("test_envs.env"); err != nil {
		logger.Error("Error loading environment variables, shutting down the app", err)
		panic(err)
	}
	Config = config{
		"users_postgres_username": os.Getenv("users_postgres_username"),
		"users_postgres_password": os.Getenv("users_postgres_password"),
		"users_postgres_schema":   os.Getenv("users_postgres_schema"),
		"oauth":                   os.Getenv("users_ouath_direction"),
		"address":                 os.Getenv("address"),
	}
}
