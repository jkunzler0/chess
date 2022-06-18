package main

import (
	"flag"

	"github.com/jkunzler0/chess/client/p2p"
)

type config struct {
	p2p       bool
	p2pConfig p2p.P2pConfig
}

func parseFlags() *config {
	c := &config{}
	flag.BoolVar(&c.p2p, "p2p", false, "P2P\n")
	flag.StringVar(&c.p2pConfig.GroupID, "group", "01", "Group ID for finding specific games\n")
	flag.StringVar(&c.p2pConfig.ListenHost, "host", "0.0.0.0", "Host listen address\n")
	flag.StringVar(&c.p2pConfig.ProtocolID, "pid", "/chess/1.0.0", "Protocol ID for stream headers\n")
	flag.IntVar(&c.p2pConfig.ListenPort, "port", 4001, "Node listen port\n")

	flag.Parse()
	return c
}
