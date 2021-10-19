package config

import (
	"os"
)

type config map[string]string

var Config config

func init() {

	Config = config{
		"database": os.Getenv("DATABASE_URL"),
		"address":  os.Getenv("address"),
		"port":     os.Getenv("PORT"),
		"oauth":    os.Getenv("oauth"),
	}
}
