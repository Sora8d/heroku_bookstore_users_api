package config

import (
	"os"
)

type config map[string]string

var Config config

func init() {

	Config = config{
		"database": os.Getenv("database"),
		"address":  os.Getenv("address"),
	}
}
