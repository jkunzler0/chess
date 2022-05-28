package main

import (
	"bufio"
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

	for {
		move, err := readMove()
		if err == nil {
			if move == "quit" || move == "q" {
				return false
			}
			err = makeMove(b, move, color)
			if err == nil {
				printBoard(*b)
				return true
			}
		}
		fmt.Println(err)
	}
}

func main() {

	fmt.Println("Chess (\"q\" or \"quit\" to quit)")

	var board Board
	err := defaultBoard(&board)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	playing := true

	for playing {

		playing = turn(&board, White)

		if !playing {
			break
		}

		playing = turn(&board, Black)
	}

	fmt.Println("Game End")
}
