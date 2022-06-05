package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func readMove() (string, error) {
	fmt.Print("Enter move: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return input, fmt.Errorf("cannot read move: %w", err)
	}

	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\r\n")
	return input, nil
}

func turn(b *Board, color bool) bool {

	var err error
	var move string
	var check [2]bool

	for {
		// Read Move
		move, err = readMove()
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Please input a valid move:")
			continue
		}
		if move == "quit" || move == "q" {
			return false
		}
		// Verify and Make Move
		err = makeMove(b, move, color)
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Please input a valid move:")
			continue
		}
		// Report Check/Checkmate and if Game is Complete
		check, err = inCheck(*b)
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Please input a valid move:")
			continue
		}
		if check[0] {
			if inCheckmate(*b, White) {
				printBoard(*b)
				fmt.Println("White is in checkmate!")
				fmt.Println("Black wins!")
				return false
			}
			fmt.Println("White is in check!")
		} else if check[1] {
			if inCheckmate(*b, Black) {
				printBoard(*b)
				fmt.Println("Black is in checkmate!")
				fmt.Println("White wins!")
				return false
			}
			fmt.Println("Black is in check!")
		}
		printBoard(*b)
		return true
	}
}

func hotseatGame() {

	fmt.Println("----- Hotsteat Chess Game -----")
	fmt.Println("For a p2p game, see `./chess -help`.")
	fmt.Println("Instructions:")
	fmt.Println("Type moves using the notation, L#L#, in which L is a letter and # is a number.")
	fmt.Println("Type \"q\" or \"quit\" to quit.")

	var board Board
	err := defaultBoard(&board)
	if err != nil {
		panic(err)
	}
	printBoard(board)

	playing := true
	for playing {
		fmt.Println("White's Turn")
		playing = turn(&board, White)
		if !playing {
			break
		}
		fmt.Println("Black's Turn")
		playing = turn(&board, Black)
	}
	fmt.Println("Game End")
}

func p2pGame(rw *bufio.ReadWriter) {

	stdReader := bufio.NewReader(os.Stdin)
	go writeStream(rw, stdReader)
	go readStream(rw)

}

func main() {
	help := flag.Bool("help", false, "Display Help")
	cfg := parseFlags()

	if *help {
		fmt.Printf("Chess!\n")
		fmt.Printf("Usage:\nRun './chess' for local hotseat game\nor\nRun './chess -p2p' to connect to and play against a local peer\n")
		os.Exit(0)
	}

	if !cfg.p2p {
		hotseatGame()
		return
	}

	p2pSetup(cfg)

}
