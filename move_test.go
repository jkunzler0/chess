package main

import (
	"testing"
)

func TestMakeMove(t *testing.T) {

	// Setup
	var b Board
	var err error
	defaultBoard(&b)
	printBoard(&b)

	// Make simple move
	err = makeMove(&b, "a2 a3", White)
	if err != nil {
		t.Error(err)
	}
	printBoard(&b)

	// Series of valid moves
	err = makeMove(&b, "e2 e4", White)
	if err != nil {
		t.Error(err)
	}
	err = makeMove(&b, "g1 f3", White)
	if err != nil {
		t.Error(err)
	}
	err = makeMove(&b, "d1 e2", White)
	if err != nil {
		t.Error(err)
	}
	// err = makeMove(&b, "e1 d1", White)
	// if err != nil {
	// 	t.Error(err)
	// }
	// err = makeMove(&b, "h1 e1", White)
	// if err != nil {
	// 	t.Error(err)
	// }
	printBoard(&b)
	// Make series of valid moves

	// Make series of invalid moves
	//
}
