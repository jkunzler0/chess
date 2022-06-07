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

	if !cfg.p2p {
		game.HotseatGame()
		return
	}

	ch := make(chan *p2p.GameHello)
	err := p2p.P2pSetup(p2pCfg, ch)
	if err != nil {
		panic(err)
	}

	fmt.Println("Start Waiting")
	gh := <-ch
	fmt.Println("STOP Waiting")
	game.P2pGame(gh.Rw, gh.White)

	select {}

}
