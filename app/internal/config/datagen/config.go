package datagen

import "flag"

type Config struct {
	DSN string

	EmployeesCount int
}

func GetConfig() Config {
	cfg := Config{}
	flag.StringVar(&cfg.DSN, "dsn", "", "DB DSN")
	flag.IntVar(&cfg.EmployeesCount, "emp-count", 10000, "Number of employees entries to generate")
	flag.Parse()

	return cfg
}
