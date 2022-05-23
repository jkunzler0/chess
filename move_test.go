package main

import (
	"testing"
)

func TestMakeMove(t *testing.T) {
	// TODO actually make this test robust

	// Setup
	var b Board
	var err error
	defaultBoard(&b)
	printBoard(&b)

	// Series of valid moves
	validMoves := []string{"a2a3", "e2 e4", "g1 f3", "d1 e2", "e2 b5", "e1d1", "f1 e2", "h1 e1"}
	for _, x := range validMoves {
		err = makeMove(&b, x, White)
		if err != nil {
			t.Error(x, err)
		}
	}
	printBoard(&b)

	// Series of invalid moves
	invalidMoves := []string{"d4d5", "i2a3", "e4e4", "g8h6", "b5e2", "b1a3", "a3a2", "e2 e3", "e1 e3", "d1f1", "f3g4"}
	for _, x := range invalidMoves {
		err = makeMove(&b, x, White)
		if err == nil {
			t.Error(x, " expects error")
		}
	}

	// Test black and pawn movement
	// moves := []string{""}
	// for _, x := range moves {
	// 	err = makeMove(&b, x, Black)
	// 	if err == nil {
	// 		t.Error(x, " expects error")
	// 	}
	// }

}
