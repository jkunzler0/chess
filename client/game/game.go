package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// #######################################################################
// (Section 1) GameState #################################################
// #######################################################################

type GameState struct {
	brd       *board
	whiteTurn bool
	reader    *bufio.Reader
	rch       chan string
	wch       chan string
}

type P2PParams struct {
	YouStart  bool
	ReadChan  chan string
	WriteChan chan string
}

func InitHotseat() (*GameState, error) {
	// Create our game broad
	b, err := defaultBoard()
	if err != nil {
		panic(err)
	}
	return &GameState{
		brd:       b,
		whiteTurn: true, // White always starts first
		reader:    bufio.NewReader(os.Stdin),
	}, nil
}

func InitP2P(p P2PParams) (*GameState, error) {
	// Create our game broad
	b, err := defaultBoard()
	if err != nil {
		panic(err)
	}
	return &GameState{
		brd:       b,
		whiteTurn: p.YouStart,
		reader:    bufio.NewReader(os.Stdin),
		rch:       p.ReadChan,
		wch:       p.WriteChan,
	}, nil
}

// #######################################################################
// (Section 2) Turns #####################################################
// #######################################################################

func readYourMove(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter move: ")
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return input, fmt.Errorf("cannot read move: %w", err)
	}

	// Remove the delimeter from the string
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	// fmt.Println("input: (", input, ")")
	return input, nil
}

func (gs *GameState) yourTurn() (bool, string) {

	var err error
	var move string
	var checkmate bool

	for {
		// Read Move
		move, err = readYourMove(gs.reader)
		if err != nil {
			fmt.Println("Error: ", err, "Move: ", move)
			fmt.Println("Please input a valid move:")
			continue
		}
		if move == "quit" || move == "q" {
			return false, move
		}
		// Verify and Make Move
		err = makeMove(gs.brd, move, gs.whiteTurn)
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Please input a valid move:")
			continue
		}
		gs.brd.printBoard()
		// Report Check/Checkmate and if Game is Complete
		checkmate, err = reportCheckAndCheckmate(*gs.brd)
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Please input a valid move:")
			continue
		}
		return !checkmate, move
	}
}

func (gs *GameState) theirTurn(move string) bool {

	var checkmate bool

	if move == "quit" || move == "q" {
		return false
	}
	// Verify and Make Move
	err := makeMove(gs.brd, move, !gs.whiteTurn)
	if err != nil {
		fmt.Println("They gave you a bad input... (", move, ")")
		panic(err)
	}
	gs.brd.printBoard()
	// Report Check/Checkmate and if Game is Complete
	checkmate, err = reportCheckAndCheckmate(*gs.brd)
	if err != nil {
		fmt.Println("They gave you a bad input... (", move, ")")
		panic(err)
	}
	return !checkmate
}

// #######################################################################
// (Section 3) Main Game Functions/Loops #################################
// #######################################################################

func (gs *GameState) PlayHotseat() {

	fmt.Println("----- Hotsteat Chess Game -----")
	fmt.Println("For a p2p game or game instructions, see `./chess -help`.")
	gs.brd.printBoard()

	playing := true
	for playing {
		if gs.whiteTurn {
			fmt.Println("White's Turn")
		} else {
			fmt.Println("Black's Turn")
		}
		playing, _ = gs.yourTurn()
		gs.whiteTurn = !gs.whiteTurn
	}
	fmt.Println("Game End")
}

func (gs *GameState) PlayP2P() (bool, bool) {

	defer close(gs.rch)
	defer close(gs.wch)

	fmt.Println("----- P2P Chess Game -----")
	fmt.Println("For a hotseat game or game instructions, see `./chess -help`.")
	gs.brd.printBoard()

	var move string
	turn := gs.whiteTurn
	playing := true
	for playing {
		if turn {
			fmt.Println("Your Turn")
			// Make your turn locally
			playing, move = gs.yourTurn()
			// Send your move to your opponent
			gs.wch <- move
		} else {
			fmt.Println("Opponents Turn")
			// Block until your opponent sends their move
			move = <-gs.rch
			// Make your opponent's move locally
			playing = gs.theirTurn(move)
			fmt.Println("Their move: ", move)
		}
		turn = !turn
	}

	fmt.Println("Game End")

	if !turn {
		fmt.Println("~~~You Win!~~~")
	} else {
		fmt.Println("~~~You Lose~~~")
	}

	if move == "quit" || move == "q" {
		return false, false
	}

	return true, !turn
}
