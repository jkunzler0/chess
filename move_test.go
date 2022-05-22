package main

import "testing"

func TestMakeMove(t *testing.T) {
	var b Board
	defaultBoard(&b)
	printBoard(&b)
	makeMove(&b, "a2 a3", true)
	printBoard(&b)
}
