package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jkunzler0/chess/client/game"
	"github.com/jkunzler0/chess/client/p2p"
	"github.com/jkunzler0/chess/client/report"
)

func main() {
	var err error
	help := flag.Bool("help", false, "Display Help")
	cfg := parseFlags()

	if *help {
		fmt.Printf("Chess!\nUsage:\nRun './chess' for local hotseat game\nor\nRun './chess -p2p' to connect to and play against a local peer\n")
		fmt.Printf("Game Instructions:\nType moves using the notation, L#L#, in which L is a letter and # is a number.")
		fmt.Println("Type \"q\" or \"quit\" to quit.")
		os.Exit(0)
	}

	var g *game.GameState

	// If p2p is off, start a hotseat game
	if !cfg.p2p {
		// Initialize GameState
		g, err = game.InitHotseat()
		if err != nil {
			panic(err)
		}
		g.PlayHotseat()
		return
	}

	// Setup p2p: providing its config and a channel to recieve the GameHello
	// The GameHello contains two channels for reading/writing to/from a peer
	//		and the color of this player
	ch := make(chan *p2p.GameHello, 1)
	err = p2p.P2pSetup(&cfg.p2pConfig, ch)
	if err != nil {
		panic(err)
	}

	// Block here until we connect to a peer
	// On connection to a peer, we receive the GameHello on ch
	gh := <-ch

	// Create the GameState with the GameHello's information
	g, err = game.InitP2P(game.P2PParams{
		YouStart:  gh.White,
		ReadChan:  gh.RCh,
		WriteChan: gh.WCh})
	if err != nil {
		panic(err)
	}

	// Exchange names with the peer
	gh.WCh <- cfg.nickname
	peerNickname := <-gh.RCh
	fmt.Printf("Connected to %s\n", peerNickname)

	// Start the P2P game
	complete, win := g.PlayP2P()
	if !complete {
		fmt.Println("Game ended in a draw.")
		return
	}

	// Report the result of the game to the server
	report.ReportResult(cfg.nickname, peerNickname, win)

}
