package employees

import (
	"os"
)

type Config struct {
	Host string
	DSN  string
}

func GetConfig() Config {
	cfg := Config{}

	var ok bool
	cfg.Host, ok = os.LookupEnv("EMPLOYEES_HOST")
	if !ok {
		cfg.Host = "localhost:8080"
	}
	cfg.DSN, _ = os.LookupEnv("EMPLOYEES_DSN")

	return cfg
}
