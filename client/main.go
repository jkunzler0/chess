package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jkunzler0/chess/client/game"
	"github.com/jkunzler0/chess/client/p2p"
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

	// Initialize Game State
	gs := game.NewGameState()

	// If p2p is off, start a hotseat game
	if !cfg.p2p {
		game.HotseatGame(gs)
		return
	}

	// Setup p2p: providing its config and a channel to recieve the GameHello
	// The GameHello contains two channels for reading/writing to/from a peer
	//		and the color of this player
	ch := make(chan *p2p.GameHello, 1)
	err := p2p.P2pSetup(p2pCfg, ch)
	if err != nil {
		panic(err)
	}

	// Block here until we connect to a peer
	// On connection to a peer, we receive the GameHello on ch
	gh := <-ch
	defer close(gh.RCh)
	defer close(gh.WCh)

	// Update the game state with the GameHello's information
	game.UpdateGameState(gs, gh.White, gh.RCh, gh.WCh)
	// Start the P2P game
	game.P2pGame(gs)

}
