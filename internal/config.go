package internal

import "flag"

type Config struct {
	Host string
	Port int
}

func ReadConfig() Config {
	var cfg Config

	flag.StringVar(&cfg.Host, "host", "localhost", "add host address for server starting")
	flag.IntVar(&cfg.Port, "port", 8080, "add port where server starting")

	flag.Parse()

	return cfg
}
