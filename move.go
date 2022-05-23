package main

import (
	"fmt"
	"strings"
	"unicode"
)

type Move struct {
	X1, Y1, X2, Y2 int
	StartPiece     rune
	EndPiece       rune
	Color          bool
	Board          *Board
}

func makeMove(b *Board, s string, color bool) error {
	// s := "a2 a3"
	// pos := "(0, 8-2), (0, 8-2)

	pos := strings.ReplaceAll(s, " ", "")
	if len(pos) != 4 {
		return fmt.Errorf("invalid bounds for move")
	}

	x1, y1, x2, y2 := int(8-(pos[1]-48)), int(pos[0]-97), int(8-(pos[3]-48)), int(pos[2]-97)

	// Check bounds on the move
	lessThanEight := x1 < 8 && x2 < 8 && y1 < 8 && y2 < 8
	atLeastZero := 0 <= x1 && 0 <= x2 && 0 <= y1 && 0 <= y2
	if !lessThanEight || !atLeastZero {
		return fmt.Errorf("invalid bounds for move")
	}

	m := Move{x1, y1, x2, y2, b[x1][y1], b[x2][y2], color, b}

	if validateMove(m) {
		m.Board[m.X1][m.Y1], m.Board[m.X2][m.Y2] = '-', m.Board[m.X1][m.Y1]
	}

	return nil

}

// const N, NE, E, SE, S, SW, W, NW = {0,1}, {1,1}
// TODO cleanup directions, with const?
// type Piece rune
// TODO replace StartPiece and EndPiece with type Piece pointers

var directions = map[rune][][2]int{'P': {{0, 1}, {0, 2}, {1, 1}, {-1, 1}},
	'N': {{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}},
	'B': {{1, -1}, {-1, 1}, {-1, 0}, {-1, -1}},
	'R': {{0, 1}, {0, -1}, {1, 0}, {-1, 0}},
	'Q': {{0, 1}, {0, -1}, {1, 1}, {1, 0}, {1, -1}, {-1, 1}, {-1, 0}, {-1, -1}},
	'K': {{0, 1}, {0, -1}, {1, 1}, {1, 0}, {1, -1}, {-1, 1}, {-1, 0}, {-1, -1}}}

func validateMove(m Move) bool {

	// Return if start is not my color OR if end is my color
	validStart := strings.Contains("prnbqkPRNBQK", string(m.StartPiece)) && (m.Color && unicode.IsUpper(m.StartPiece) || !m.Color && unicode.IsLower(m.StartPiece))
	validEnd := m.EndPiece == '-' || m.Color && unicode.IsLower(m.EndPiece) || !m.Color && unicode.IsUpper(m.EndPiece)
	if !validStart || !validEnd {
		return false
	}

	// TODO cleanup this switch?
	switch m.StartPiece {
	case 'P', 'p':
		validateMoveJump(m)
		// TODO all the rules for pawn
	case 'N', 'n':
		validateMoveJump(m)
	case 'K', 'k':
		validateMoveJump(m)
		// TODO Restrict movement into checkmate
	case 'Q', 'B', 'q', 'b':
		validateMoveCrawl(m)
	case 'R', 'r':
		validateMoveCrawl(m)
		// TODO Add castling
	}

	return true
}

func validateMoveJump(m Move) bool {
	j := [2]int{m.X2 - m.X1, m.Y2 - m.Y1}
	for _, i := range directions[m.StartPiece] {
		if i == j {
			return true
		}
	}
	return false
}
func validateMoveCrawl(m Move) bool {

	return true
}

// func flipMove(s string) {
// 	// makeMove()
// }
