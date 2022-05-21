package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var playing bool
var move string

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
	input = strings.TrimSuffix(input, "\n")
	// fmt.Println(input)
	return input, nil
}

func turn(p int) {
	move, _ := readMove()

	fmt.Println(move)
	// if p == 1 {
	// 	makeMove(board1, move)
	// 	makeMove(board2, flipMove(move))
	// } else {
	// 	makeMove(board2, move)
	// 	makeMove(board1, flipMove(move))
	// }
}

func main() {

	// board1 = newBoard()
	// board2 = newBoard()

	// playing = true

	// for playing {
	// 	turn(1)

	// 	if !playing {
	// 		break
	// 	}

	// 	turn(2)
	// }

}
