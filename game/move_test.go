package game

import (
	"testing"
)

func TestMakeMove(t *testing.T) {
	// TODO actually make this test robust

	// Setup
	var b Board
	var err error
	defaultBoard(&b)
	// printBoard(b)

	// Series of valid moves
	validMoves := []string{"a2a3", "e2 e4", "g1 f3", "d1 e2", "e2 b5", "e1d1", "f1 e2", "h1 e1"}
	for _, x := range validMoves {
		err = makeMove(&b, x, White)
		if err != nil {
			t.Error(x, err)
			break
		}
	}
	// printBoard(b)

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
	// printBoard(b)

}

func TestInCheck(t *testing.T) {
	var b Board
	var err error
	var check [2]bool
	var checkmate bool

	// White in check, black NOT in check
	err = newBoard(&b, "5R2/8/4k3/8/8/r2K4/8/8")
	if err != nil {
		t.Error(err)
	}
	// printBoard(b)
	check, err = inCheck(b)
	if check != [2]bool{true, false} {
		t.Error(check, " failed inCheck: ", err)
	}
	checkmate = inCheckmate(b, White)
	if checkmate != false {
		t.Error(checkmate, " failed inCheckmate: ")
	}
	// printBoard(b)

	// White in checkmate, black in check
	err = newBoard(&b, "4R3/6r1/3k4/8/8/5r1K/8/7q")
	if err != nil {
		t.Error(err)
	}
	// printBoard(b)
	check, err = inCheck(b)
	if check != [2]bool{true, false} {
		t.Error(check, " failed inCheck: ", err)
	}
	checkmate = inCheckmate(b, White)
	if checkmate != true {
		t.Error(checkmate, " failed inCheckmate: ")
	}
	// printBoard(b)

	// White NOT in check, black in check
	err = newBoard(&b, "4r2R/8/3K4/8/8/7k/8/6Q1")
	if err != nil {
		t.Error(err)
	}
	// printBoard(b)
	check, err = inCheck(b)
	if check != [2]bool{false, true} {
		t.Error(check, " failed inCheck: ", err)
	}
	checkmate = inCheckmate(b, Black)
	if checkmate != false {
		t.Error(checkmate, " failed inCheckmate: ")
	}
	// printBoard(b)
}

func TestMoveIntoCheck(t *testing.T) {
	// Setup
	var b Board
	var err error
	newBoard(&b, "4K3/8/8/8/3q4/8/8/3k4")
	// printBoard(b)

	// Series of invalid moves
	invalidMoves := []string{"e8d8", "e8d7", "e8e6"}
	for _, x := range invalidMoves {
		err = makeMove(&b, x, White)
		if err == nil {
			t.Error(x, " expects error")
			break
		}
		// t.Log(err)
	}

	// Series of valid moves
	validMoves := []string{"e8e7"}
	for _, x := range validMoves {
		err = makeMove(&b, x, White)
		if err != nil {
			t.Error(x, err)
			break
		}
		// t.Log(err)
	}

	// printBoard(b)
}
