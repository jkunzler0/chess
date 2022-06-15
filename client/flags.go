package main

import (
	"flag"

	"github.com/jkunzler0/chess/client/p2p"
)

type config struct {
	p2p bool
	// groupID    string
	// protocolID string
	// listenHost string
	// listenPort int
}

func parseFlags() (*config, *p2p.P2pConfig) {
	c := &config{}
	flag.BoolVar(&c.p2p, "p2p", false, "P2P\n")

	p := &p2p.P2pConfig{}
	flag.StringVar(&p.GroupID, "group", "01", "Group ID for finding specific games\n")
	flag.StringVar(&p.ListenHost, "host", "0.0.0.0", "Host listen address\n")
	flag.StringVar(&p.ProtocolID, "pid", "/chess/1.0.0", "Protocol ID for stream headers\n")
	flag.IntVar(&p.ListenPort, "port", 4001, "Node listen port\n")

	flag.Parse()
	return c, p
}
