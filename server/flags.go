package main

import (
	"flag"

	"github.com/jkunzler0/chess/server/database"
)

type config struct {
	localLog   bool
	listenPort int
	dbParams   database.PostgresDbParams
}

func parseFlags() *config {
	c := &config{}
	flag.BoolVar(&c.localLog, "ll", true, "Transaction Log Location\n")
	flag.IntVar(&c.listenPort, "port", 5000, "Server listen port\n")

	flag.StringVar(&c.dbParams.DBName, "dbname", "chess", "Database name\n")
	flag.StringVar(&c.dbParams.Host, "host", "localhost", "Database host\n")
	flag.StringVar(&c.dbParams.User, "user", "test", "Database user\n")
	flag.StringVar(&c.dbParams.Password, "password", "gg", "Database password\n")

	flag.Parse()
	return c
}
