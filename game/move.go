package game

import (
	"fmt"
	"strings"
	"unicode"
)

const White bool = true
const Black bool = false

// type Piece rune
// TODO replace StartPiece and EndPiece with type Piece pointers?

type Move struct {
	X1, Y1, X2, Y2 int
	StartPiece     rune
	EndPiece       rune
	Color          bool
	Board          Board
}

// #######################################################################
// (Section 1) Moving and Verifying Moves ################################
// #######################################################################

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
	m := Move{x1, y1, x2, y2, b[x1][y1], b[x2][y2], color, *b}
	if ok, err := validateMove(m); ok {
		b[m.X1][m.Y1], b[m.X2][m.Y2] = '-', b[m.X1][m.Y1]
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
	var valid bool
	var err error
	switch m.StartPiece {
	case 'P', 'p':

		// Extra rules for pawn
		if m.Y1-m.Y2 == 2 && (m.Y1 != 1 && m.Y1 != 6) {
			return false, fmt.Errorf("pawn can only advance 2 squares if it has not already moved")
		} else if (m.X1-m.X2 != 0 && m.EndPiece == '-') || (m.X1-m.X2 == 0 && m.EndPiece != '-') {
			return false, fmt.Errorf("pawn can only attack diagonally and move vertically")
		} else if m.Y1-m.Y2 >= 1 && !m.Color || m.Y1-m.Y2 <= -1 && m.Color {
			return false, fmt.Errorf("pawn is going the wrong way")
		}

		valid, err = validateMoveJump(m)
		// TODO add pawn promotion
	case 'N', 'n':
		valid, err = validateMoveJump(m)
	case 'K', 'k':
		valid, err = validateMoveJump(m)
	case 'Q', 'B', 'q', 'b':
		valid, err = validateMoveCrawl(m)
	case 'R', 'r':
		valid, err = validateMoveCrawl(m)
		// TODO Add castling
	default:
		return false, fmt.Errorf("not a valid piece type")
	}
	if !valid {
		return valid, err
	}

	// Verify that this move does NOT put our king into check
	m.Board[m.X1][m.Y1], m.Board[m.X2][m.Y2] = '-', m.Board[m.X1][m.Y1]
	check, _ := inCheck(m.Board)
	if unicode.IsUpper(m.StartPiece) && check[0] {
		return !check[0], fmt.Errorf("cannot put your own king into check")
	} else if unicode.IsLower(m.StartPiece) && check[1] {
		return !check[1], fmt.Errorf("cannot put your own king into check")
	}

	return true, nil
}

// Pawns, Knights, and Kings jump to a location (as opposed to crawling/sliding)
func validateMoveJump(m Move) (bool, error) {
	// fmt.Println("coor", m.X1, " ", m.Y1, " ", m.X2, " ", m.Y2)
	j := [2]int{m.X1 - m.X2, m.Y1 - m.Y2}
	for _, i := range getDirections(m.StartPiece) {
		// fmt.Println("match", i, " ", j)
		if i == j {
			return true, nil
		}
	}
	return false, fmt.Errorf("%s cannot move there", string(m.StartPiece))
}

// Queen, Bishops, and Rooks crawl/slide across the board
func validateMoveCrawl(m Move) (bool, error) {
	for _, i := range getDirections(m.StartPiece) {
		x, y := m.X1+i[0], m.Y1+i[1]
		for inBounds(x, y) {
			// fmt.Println(x, y, m.X1, m.Y1, m.X2, m.Y2, i)
			if x == m.X2 && y == m.Y2 {
				// We made it to the endPiece
				return true, nil
			} else if m.Board[x][y] != '-' {
				// We ran into another piece
				break
			}
			x += i[0]
			y += i[1]
		}
	}
	return false, fmt.Errorf("%s cannot move there", string(m.StartPiece))
}

// #######################################################################
// (Section 2) Check and Checkmate #######################################
// #######################################################################

func inCheck(b Board) ([2]bool, error) {

	// Get location of both kings
	wk, bk, err := findKings(b)
	if err != nil {
		return [2]bool{false, false}, fmt.Errorf("%w", err)
	}

	// Create moves against the kings
	mbk := Move{X2: bk[0], Y2: bk[1], Color: White, EndPiece: 'k', Board: b}
	mwk := Move{X2: wk[0], Y2: wk[1], Color: Black, EndPiece: 'K', Board: b}

	// To determine if a kings is in check,
	// attempt to validate moves of every pieces against the enemy king
	var whiteCheck, blackCheck bool
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if unicode.IsLetter(b[x][y]) {
				if unicode.IsUpper(b[x][y]) && !blackCheck {
					mbk.X1, mbk.Y1, mbk.StartPiece = x, y, b[x][y]
					blackCheck, _ = validateMove(mbk)
				} else if unicode.IsLower(b[x][y]) && !whiteCheck {
					mwk.X1, mwk.Y1, mwk.StartPiece = x, y, b[x][y]
					whiteCheck, _ = validateMove(mwk)
				}
			}
		}
	}

	return [2]bool{whiteCheck, blackCheck}, nil
}

func inCheckmate(b Board, kingColor bool) bool {

	var m Move
	tmpB := b

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if unicode.IsLetter(b[x][y]) &&
				(kingColor && unicode.IsUpper(b[x][y]) || !kingColor && unicode.IsLower(b[x][y])) {

				m = Move{X1: x, Y1: y, Color: kingColor, StartPiece: b[x][y], Board: tmpB}
				for z := 0; z < 8; z++ {
					for w := 0; w < 8; w++ {
						m.X2, m.Y2, m.EndPiece = z, w, b[z][w]
						validMove, _ := validateMove(m)
						if validMove {
							tmpB[m.X1][m.Y1], tmpB[m.X2][m.Y2] = '-', tmpB[m.X1][m.Y1]
							inCheck, _ := inCheck(tmpB)
							if kingColor && !inCheck[0] {
								return false
							} else if !kingColor && !inCheck[1] {
								return false
							} else {
								// Reset temp board
								tmpB = b
							}
						}
					}
				}
			}
		}
	}

	return true
}

// Prints if any check or checkmate
// Return true if any checkmate
func reportCheckAndCheckmate(b Board) (bool, error) {

	check, err := inCheck(b)
	if err != nil {
		return false, err
	}
	if check[0] {
		if inCheckmate(b, White) {
			printBoard(b)
			fmt.Println("White is in checkmate!")
			fmt.Println("Black wins!")
			return true, nil
		}
		fmt.Println("White is in check!")
	} else if check[1] {
		if inCheckmate(b, Black) {
			printBoard(b)
			fmt.Println("Black is in checkmate!")
			fmt.Println("White wins!")
			return true, nil
		}
		fmt.Println("Black is in check!")
	}
	return false, nil
}

// #######################################################################
// (Section 3) Helper Functions ##########################################
// #######################################################################

func inBounds(x int, y int) bool {
	if x < 8 && y < 8 && x >= 0 && y >= 0 {
		return true
	}
	return false
}

func getDirections(piece rune) [][2]int {
	var directions = map[rune][][2]int{'P': {{0, 1}, {0, 2}, {1, 1}, {-1, 1}, {0, -1}, {0, -2}, {1, -1}, {-1, -1}},
		'N': {{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}},
		'B': {{1, -1}, {-1, 1}, {-1, 0}, {-1, -1}},
		'R': {{0, 1}, {0, -1}, {1, 0}, {-1, 0}},
		'Q': {{0, 1}, {0, -1}, {1, 1}, {1, 0}, {1, -1}, {-1, 1}, {-1, 0}, {-1, -1}},
		'K': {{0, 1}, {0, -1}, {1, 1}, {1, 0}, {1, -1}, {-1, 1}, {-1, 0}, {-1, -1}}}

	if unicode.IsLetter(piece) {
		return directions[unicode.ToUpper(piece)]
	} else {
		return [][2]int{} // TODO Could return error here
	}
}

func findKings(b Board) ([2]int, [2]int, error) {

	var wk, bk [2]int
	var wkFound, bkFound bool

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if b[x][y] == 'K' {
				wk = [2]int{x, y}
				wkFound = true
			} else if b[x][y] == 'k' {
				bk = [2]int{x, y}
				bkFound = true
			}
		}
	}
	if !wkFound || !bkFound {
		return wk, bk, fmt.Errorf("board is missing kings")
	}
	return wk, bk, nil
}
