package main

import (
	"fmt"
	"strings"
	"unicode"
)

type Move struct {
	X1, Y1, X2, Y2 int
	Color          bool
}

func makeMove(b *Board, s string, color bool) error {
	// s := "a2 a3"
	// pos := "(0, 8-2), (0, 8-2)

	pos := strings.ReplaceAll(s, " ", "")
	if len(pos) != 4 {
		return fmt.Errorf("invalid board")
	}

	move := Move{int(8 - (pos[1] - 48)), int(pos[0] - 97),
		int(8 - (pos[3] - 48)), int(pos[2] - 97),
		color}

	if validateMove(b, move) {
		b[move.X1][move.Y1], b[move.X2][move.Y2] = '-', b[move.X1][move.Y1]
	}

	return nil

}

var ability = map[rune]int{'r': 2, h: 3}

func validateMove(b *Board, move Move) bool {

	// if in bounds

	start := b[move.X1][move.Y1]
	end := b[move.X2][move.Y2]

	// if start is not my piece or end is my piece, return
	validStart := move.Color && unicode.IsUpper(start) || !move.Color && unicode.IsLower(start)
	validEnd := end == '-' || move.Color && unicode.IsLower(end) || !move.Color && unicode.IsUpper(start)
	if !validStart || !validEnd {
		return false
	}

	// if have the ability to move there
	// if there is nothing in my way
	return true
}

// func flipMove(s string) {
// 	// makeMove()
// }
