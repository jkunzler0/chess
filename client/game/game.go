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

type gameState struct {
	brd       *board
	whiteTurn bool
	reader    *bufio.Reader
	rch       chan string
	wch       chan string
}

func NewGameState() *gameState {
	// Create our game broad
	var b board
	err := defaultBoard(&b)
	if err != nil {
		panic(err)
	}

	return &gameState{
		brd:       &b,
		whiteTurn: true, // White always starts first
		reader:    bufio.NewReader(os.Stdin),
	}
}

// UpdateGameState to use channels to communicate moves for p2p
func UpdateGameState(gs *gameState, whiteTurn bool, rch chan string, wch chan string) error {
	gs.whiteTurn = whiteTurn
	gs.rch = rch
	gs.wch = wch

	return nil
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

func yourTurn(gs *gameState) (bool, string) {

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
		printBoard(*gs.brd)
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

func theirTurn(b *board, color bool, move string) bool {

	var checkmate bool

	if move == "quit" || move == "q" {
		return false
	}
	// Verify and Make Move
	err := makeMove(b, move, color)
	if err != nil {
		fmt.Println("They gave you a bad input... (", move, ")")
		panic(err)
	}
	printBoard(*b)
	// Report Check/Checkmate and if Game is Complete
	checkmate, err = reportCheckAndCheckmate(*b)
	if err != nil {
		fmt.Println("They gave you a bad input... (", move, ")")
		panic(err)
	}
	return !checkmate
}

// #######################################################################
// (Section 3) Main Game Functions/Loops #################################
// #######################################################################

func HotseatGame(gs *gameState) {

	fmt.Println("----- Hotsteat Chess Game -----")
	fmt.Println("For a p2p game or game instructions, see `./chess -help`.")
	printBoard(*gs.brd)

	playing := true
	for playing {
		if gs.whiteTurn {
			fmt.Println("White's Turn")
		} else {
			fmt.Println("Black's Turn")
		}
		playing, _ = yourTurn(gs)
		gs.whiteTurn = !gs.whiteTurn
	}
	fmt.Println("Game End")
}

func P2pGame(gs *gameState) {

	fmt.Println("----- P2P Chess Game -----")
	fmt.Println("For a hotseat game or game instructions, see `./chess -help`.")
	printBoard(*gs.brd)

	var move string
	turn := gs.whiteTurn
	playing := true
	for playing {
		if turn {
			fmt.Println("Your Turn")
			// Make your turn locally
			playing, move = yourTurn(gs)
			// Send your move to your opponent
			gs.wch <- move
		} else {
			fmt.Println("Opponents Turn")
			// Block until your opponent sends their move
			move = <-gs.rch
			// Make your opponent's move locally
			playing = theirTurn(gs.brd, !gs.whiteTurn, move)
			fmt.Println("Their move: ", move)
		}
		turn = !turn
	}
	fmt.Println("Game End")
}
