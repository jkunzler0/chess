package main

import (
	"flag"
)

type config struct {
	p2p        bool
	groupID    string
	protocolID string
	listenHost string
	listenPort int
}

func parseFlags() *config {
	c := &config{}

	flag.BoolVar(&c.p2p, "p2p", false, "P2P\n")
	flag.StringVar(&c.groupID, "group", "01", "Group ID for finding specific games\n")
	flag.StringVar(&c.listenHost, "host", "0.0.0.0", "Host listen address\n")
	flag.StringVar(&c.protocolID, "pid", "/chess/1.0.0", "Protocol ID for stream headers\n")
	flag.IntVar(&c.listenPort, "port", 4001, "Node listen port\n")

	flag.Parse()
	return c
}
