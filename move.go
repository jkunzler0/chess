package main

import (
	"fmt"
	"strings"
	"unicode"
)

const White bool = true
const Black bool = false

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
	if !inBounds(x1, y1) || !inBounds(x2, y2) {
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
		return validateMoveJump(m)
		// TODO all the rules for pawn
	case 'N', 'n':
		return validateMoveJump(m)
	case 'K', 'k':
		return validateMoveJump(m)
		// TODO Restrict movement into checkmate
	case 'Q', 'B', 'q', 'b':
		return validateMoveCrawl(m)
	case 'R', 'r':
		return validateMoveCrawl(m)
		// TODO Add castling
	}

	return false
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
	for _, i := range directions[m.StartPiece] {
		x, y := m.X1, m.Y1
		for inBounds(x, y) {
			x += i[0]
			y += i[1]

			// We made it to the endPiece
			if x == m.X2 && y == m.Y2 {
				return true
			} else if m.Board[x][y] != '-' {
				break
			}
		}
	}
	return false
}

func inBounds(x int, y int) bool {
	if x < 8 && y < 8 && x >= 0 && y >= 0 {
		return true
	}
	return false
}

// func flipMove(s string) {
// 	// makeMove()
// }
