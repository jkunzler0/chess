package main

import (
	"flag"

	"github.com/jkunzler0/chess/server/database"
)

type config struct {
	postgresLog  bool
	listenPort   int
	localLogName string
	dbParams     database.PostgresDbParams
}

func parseFlags() *config {
	c := &config{}
	flag.BoolVar(&c.postgresLog, "pg", false, "Transaction Log Location\n")
	flag.IntVar(&c.listenPort, "port", 8080, "Server listen port\n")
	flag.StringVar(&c.localLogName, "ll", "tmp", "	Local Log Name\n")

	flag.StringVar(&c.dbParams.DBName, "dbname", "chess", "Database name\n")
	flag.StringVar(&c.dbParams.Host, "host", "localhost", "Database host\n")
	flag.StringVar(&c.dbParams.User, "user", "test", "Database user\n")
	flag.StringVar(&c.dbParams.Password, "password", "gg", "Database password\n")

	flag.Parse()
	return c
}
