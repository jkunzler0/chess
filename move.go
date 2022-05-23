package main

import (
	"fmt"
	"strings"
	"unicode"
)

const White bool = true
const Black bool = false

// type Piece rune
// TODO replace StartPiece and EndPiece with type Piece pointers

type Move struct {
	X1, Y1, X2, Y2 int
	StartPiece     rune
	EndPiece       rune
	Color          bool
	Board          *Board
}

func makeMove(b *Board, s string, color bool) error {
	// e.g.	`s := "a2 a3"` OR `s := "a2a3"`

	// Parse move string, s, into coordinates
	pos := strings.ReplaceAll(s, " ", "")
	if len(pos) != 4 {
		return fmt.Errorf("invalid move format")
	}
	x1, y1, x2, y2 := int(pos[0]-97), int(8-(pos[1]-48)), int(pos[2]-97), int(8-(pos[3]-48))

	// fmt.Println(x1, y1, x2, y2)

	// Check that the coordinates are within the bounds of the board
	if !inBounds(x1, y1) || !inBounds(x2, y2) {
		return fmt.Errorf("invalid bounds for move")
	}

	// Create a Move struct, validate the move, and then make the move
	m := Move{x1, y1, x2, y2, b[x1][y1], b[x2][y2], color, b}
	if ok, err := validateMove(m); ok {
		m.Board[m.X1][m.Y1], m.Board[m.X2][m.Y2] = '-', m.Board[m.X1][m.Y1]
	} else {
		return fmt.Errorf("move is invalid: %w", err)
	}

	return nil
}

// Return true if the move is valid, return false otherwise
func validateMove(m Move) (bool, error) {

	// Return if start is not my color OR if end is my color
	validStart := strings.Contains("prnbqkPRNBQK", string(m.StartPiece)) && (m.Color && unicode.IsUpper(m.StartPiece) || !m.Color && unicode.IsLower(m.StartPiece))
	validEnd := m.EndPiece == '-' || m.Color && unicode.IsLower(m.EndPiece) || !m.Color && unicode.IsUpper(m.EndPiece)
	if !validStart || !validEnd {
		return false, fmt.Errorf("invalid start or invalid end")
	}

	// Check rules for specific pieces
	switch m.StartPiece {
	case 'P', 'p':
		return validateMoveJump(m)
		// TODO all the extra rules for pawn
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

	return false, fmt.Errorf("not a valid piece type")
}

// Pawns, Knights, and Kings jump to a location (as opposed to crawling/sliding)
func validateMoveJump(m Move) (bool, error) {
	directions := getDirections()
	// fmt.Println(m.X1, " ", m.Y1, " ", m.X2, " ", m.Y2)
	j := [2]int{m.X1 - m.X2, m.Y1 - m.Y2}
	for _, i := range directions[m.StartPiece] {
		// fmt.Println(i, " ", j)
		if i == j {
			return true, nil
		}
	}
	return false, fmt.Errorf("failled validateMoveJump")
}

// Queen, Bishops, and Rooks crawl/slide across the board
func validateMoveCrawl(m Move) (bool, error) {
	directions := getDirections()
	for _, i := range directions[m.StartPiece] {
		x, y := m.X1, m.Y1
		for inBounds(x, y) {
			x += i[0]
			y += i[1]

			if x == m.X2 && y == m.Y2 {
				// We made it to the endPiece
				return true, nil
			} else if m.Board[x][y] != '-' {
				// We ran into another piece
				break
			}
		}
	}
	return false, fmt.Errorf("failled validateMoveCrawl")
}

// Helper functions
func inBounds(x int, y int) bool {
	if x < 8 && y < 8 && x >= 0 && y >= 0 {
		return true
	}
	return false
}
func getDirections() map[rune][][2]int {
	// Hide piece directions in a get function to avoid global variables
	return map[rune][][2]int{'P': {{0, 1}, {0, 2}, {1, 1}, {-1, 1}},
		'N': {{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}},
		'B': {{1, -1}, {-1, 1}, {-1, 0}, {-1, -1}},
		'R': {{0, 1}, {0, -1}, {1, 0}, {-1, 0}},
		'Q': {{0, 1}, {0, -1}, {1, 1}, {1, 0}, {1, -1}, {-1, 1}, {-1, 0}, {-1, -1}},
		'K': {{0, 1}, {0, -1}, {1, 1}, {1, 0}, {1, -1}, {-1, 1}, {-1, 0}, {-1, -1}}}
}
