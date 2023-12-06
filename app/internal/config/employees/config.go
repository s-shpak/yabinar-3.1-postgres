package employees

import "flag"

type Config struct {
	Host string
	DSN  string
}

func GetConfig() Config {
	cfg := Config{}

	flag.StringVar(&cfg.Host, "host", "localhost:8080", "server listen address")
	flag.StringVar(&cfg.DSN, "dsn", "", "DB DSN")
	flag.Parse()

	return cfg
}
