package main

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	var b Board
	newBoard(&b, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	// t.Log(b)
	printBoard(&b)

	var c Board
	defaultBoard(&c)
	printBoard(&c)

}
