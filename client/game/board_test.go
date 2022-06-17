package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	// var err error

	// Default board
	_, err := newBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	if err != nil {
		t.Error(err)
	}
	// b.printBoard()
	// Misc board
	_, err = newBoard("8/4pK2/8/qr6/8/8/PPPPPPPP/RNBQ3R")
	if err != nil {
		t.Error(err)
	}

	// Rewrite previous board
	_, err = defaultBoard()
	if err != nil {
		t.Error(err)
	}

	// Invalid boards
	_, err = newBoard("8/4pK2/8/8/8/PPPPPPPP/RNBQ3R")
	if err == nil {
		t.Error("expected error, missing rank/row")
	}
	_, err = newBoard("rnbqkbnr/ppppppp/8/8/8/8/PPPPPPPP/RNQKBNR")
	if err == nil {
		t.Error("expected error, missing pieces")
	}

}
