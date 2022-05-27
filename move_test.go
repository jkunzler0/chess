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
			break
		}
	}
	printBoard(&b)

	// Series of invalid moves
	invalidMoves := []string{"d4d5", "i2a3", "e4e4", "g8h6", "b5e2", "b1a3", "a3a2", "e2 e3", "e1 e3", "d1f1", "f3g4"}
	for _, x := range invalidMoves {
		err = makeMove(&b, x, White)
		if err == nil {
			t.Error(x, " expects error")
			break
		}
	}

	// Test black and pawn movement
	moves := []string{"f7f5", "f5e4"}
	for _, x := range moves {
		err = makeMove(&b, x, Black)
		if err != nil {
			t.Error(x, err)
			break
		}
	}
	printBoard(&b)

}

func TestInCheck(t *testing.T) {
	var b Board
	var err error

	// White in check, black NOT in check
	err = newBoard(&b, "5R2/8/4k3/8/8/r2K4/8/8")
	if err != nil {
		t.Error(err)
	}
	printBoard(&b)

	var check [2]bool
	check, err = inCheck(&b)
	if check != [2]bool{true, false} {
		t.Error(check, " failed inCheck: ", err)
	}

	var checkmate [2]bool
	checkmate, err = inCheckmate(&b)
	if checkmate != [2]bool{false, false} {
		t.Error(check, " failed inCheckmate: ", err)
	}

	// White in checkmate, black in check
	err = newBoard(&b, "4R3/8/4k3/8/8/r6K/8/6q1")
	if err != nil {
		t.Error(err)
	}
	printBoard(&b)
}
