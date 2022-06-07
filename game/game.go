package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readYourMove() (string, error) {
	fmt.Print("Enter move: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return input, fmt.Errorf("cannot read move: %w", err)
	}

	// Remove the delimeter from the string
	input = strings.TrimSuffix(input, "\r\n")
	return input, nil
}

func yourTurn(b *Board, color bool) (bool, string) {

	var err error
	var move string
	var checkmate bool

	for {
		// Read Move
		move, err = readYourMove()
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Please input a valid move:")
			continue
		}
		if move == "quit" || move == "q" {
			return false, move
		}
		// Verify and Make Move
		err = makeMove(b, move, color)
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Please input a valid move:")
			continue
		}
		printBoard(*b)
		// Report Check/Checkmate and if Game is Complete
		checkmate, err = reportCheckAndCheckmate(*b)
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Please input a valid move:")
			continue
		}
		return !checkmate, move
	}
}

func HotseatGame() {

	fmt.Println("----- Hotsteat Chess Game -----")
	fmt.Println("For a p2p game or game instructions, see `./chess -help`.")

	var board Board
	err := defaultBoard(&board)
	if err != nil {
		panic(err)
	}
	printBoard(board)

	playing, color := true, true
	for playing {
		if color {
			fmt.Println("White's Turn")
		} else {
			fmt.Println("Black's Turn")
		}
		playing, _ = yourTurn(&board, color)
		color = !color
	}
	fmt.Println("Game End")
}

func theirTurn(b *Board, color bool, move string) bool {

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

func P2pGame(rw *bufio.ReadWriter, yourColor bool) {

	fmt.Println("----- P2P Chess Game -----")
	fmt.Println("For a hotseat game or game instructions, see `./chess -help`.")

	var board Board
	err := defaultBoard(&board)
	if err != nil {
		panic(err)
	}
	printBoard(board)

	var move string
	turn := yourColor
	playing := true
	for playing {
		if turn {
			fmt.Println("Your Turn")
			// Make your turn locally
			playing, move = yourTurn(&board, yourColor)
			// Send your move to your opponent
			// p2p.WriteStream(rw, move)
		} else {
			fmt.Println("Opponents Turn")
			// Wait for your opponent to send their move
			// move = p2p.ReadStream(rw)
			// Make your opponent's move locally
			playing = theirTurn(&board, !yourColor, move)
			fmt.Println(move)
		}
		turn = !turn
	}
	fmt.Println("Game End")
	os.Exit(0)
}
