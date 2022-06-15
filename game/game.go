package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type gameState struct {
	brd    *board
	white  bool
	reader *bufio.Reader
	rch    chan string
	wch    chan string
}

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
			select {}
			// continue
		}
		if move == "quit" || move == "q" {
			return false, move
		}
		// Verify and Make Move
		err = makeMove(gs.brd, move, gs.white)
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

func HotseatGame(gs *gameState) {

	fmt.Println("----- Hotsteat Chess Game -----")
	fmt.Println("For a p2p game or game instructions, see `./chess -help`.")
	printBoard(*gs.brd)

	playing, color := true, true
	for playing {
		if color {
			fmt.Println("White's Turn")
		} else {
			fmt.Println("Black's Turn")
		}
		playing, _ = yourTurn(gs)
		color = !color
	}
	fmt.Println("Game End")
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

func P2pGame(gs *gameState) {

	fmt.Println("----- P2P Chess Game -----")
	fmt.Println("For a hotseat game or game instructions, see `./chess -help`.")
	printBoard(*gs.brd)

	var move string
	turn := gs.white
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
			playing = theirTurn(gs.brd, !gs.white, move)
			fmt.Println("Their move: ", move)
		}
		turn = !turn
	}
	fmt.Println("Game End")
}

func NewGameState() *gameState {

	// Create our game brd
	var b board
	err := defaultBoard(&b)
	if err != nil {
		panic(err)
	}

	return &gameState{
		brd:    &b,
		white:  true,
		reader: bufio.NewReader(os.Stdin),
	}
}

func UpdateGameState(gs *gameState, white bool, rch chan string, wch chan string) error {

	gs.white = white
	gs.rch = rch
	gs.wch = wch

	return nil
}
