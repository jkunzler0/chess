package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jkunzler0/chess/game"
	"github.com/jkunzler0/chess/p2p"
)

func main() {
	help := flag.Bool("help", false, "Display Help")
	cfg, p2pCfg := parseFlags()

	if *help {
		fmt.Printf("Chess!\nUsage:\nRun './chess' for local hotseat game\nor\nRun './chess -p2p' to connect to and play against a local peer\n")
		fmt.Printf("Game Instructions:\nType moves using the notation, L#L#, in which L is a letter and # is a number.")
		fmt.Println("Type \"q\" or \"quit\" to quit.")
		os.Exit(0)
	}

	// If p2p is off, start a hotseat game
	if !cfg.p2p {
		game.HotseatGame()
		return
	}

	// Setup p2p, passing its config and a channel to hear back from
	ch := make(chan *p2p.GameHello, 1)
	err := p2p.P2pSetup(p2pCfg, ch)
	if err != nil {
		panic(err)
	}

	// Block here until we connect to a peer
	// On connection to a peer, we receive the game hello on ch
	gh := <-ch
	game.P2pGame(gh.RCh, gh.WCh, gh.White)

}
