package main

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	var err error

	// Default board
	var b Board
	err = newBoard(&b, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	if err != nil {
		t.Error(err)
	}
	// printBoard(b)

	// Misc board
	var c Board
	err = newBoard(&c, "8/4pK2/8/qr6/8/8/PPPPPPPP/RNBQ3R")
	if err != nil {
		t.Error(err)
	}
	// printBoard(c)

	// Rewrite previous board
	err = defaultBoard(&c)
	if err != nil {
		t.Error(err)
	}
	// printBoard(c)
	// printBoardBasic(c)

	// Invalid boards
	err = newBoard(&c, "8/4pK2/8/8/8/PPPPPPPP/RNBQ3R")
	if err == nil {
		t.Error("expected error, missing rank/row")
	}
	err = newBoard(&c, "rnbqkbnr/ppppppp/8/8/8/8/PPPPPPPP/RNQKBNR")
	if err == nil {
		t.Error("expected error, missing pieces")
	}

}
