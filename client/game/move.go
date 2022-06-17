package game

import (
	"fmt"
	"strings"
	"unicode"
)

const White bool = true

type move struct {
	x1, y1, x2, y2 int
	startPiece     rune
	endPiece       rune
	white          bool
	brd            board
}

// #######################################################################
// (Section 1) Moving and Verifying Moves ################################
// #######################################################################

func makeMove(b *board, s string, white bool) error {
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

	// Create a move struct, validate the move, and then make the move
	m := move{x1, y1, x2, y2, b[x1][y1], b[x2][y2], white, *b}
	if ok, err := validateMove(m); ok {
		b[m.x1][m.y1], b[m.x2][m.y2] = '-', b[m.x1][m.y1]
	} else {
		return fmt.Errorf("move is invalid: %w", err)
	}

	return nil
}

// Return true if the move is valid, return false otherwise
func validateMove(m move) (bool, error) {

	// Return if start is not my color OR if end is my color
	validStart := strings.Contains("prnbqkPRNBQK", string(m.startPiece)) && (m.white && unicode.IsUpper(m.startPiece) || !m.white && unicode.IsLower(m.startPiece))
	validEnd := m.endPiece == '-' || m.white && unicode.IsLower(m.endPiece) || !m.white && unicode.IsUpper(m.endPiece)
	if !validStart || !validEnd {
		return false, fmt.Errorf("invalid start or invalid end")
	}

	// Check rules for specific pieces
	var valid bool
	var err error
	switch m.startPiece {
	case 'P', 'p':

		// Extra rules for pawn
		if m.y1-m.y2 == 2 && (m.y1 != 1 && m.y1 != 6) {
			return false, fmt.Errorf("pawn can only advance 2 squares if it has not already moved")
		} else if (m.x1-m.x2 != 0 && m.endPiece == '-') || (m.x1-m.x2 == 0 && m.endPiece != '-') {
			return false, fmt.Errorf("pawn can only attack diagonally and move vertically")
		} else if m.y1-m.y2 >= 1 && !m.white || m.y1-m.y2 <= -1 && m.white {
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
	m.brd[m.x1][m.y1], m.brd[m.x2][m.y2] = '-', m.brd[m.x1][m.y1]
	check, _ := inCheck(m.brd)
	if unicode.IsUpper(m.startPiece) && check[0] {
		return !check[0], fmt.Errorf("cannot put your own king into check")
	} else if unicode.IsLower(m.startPiece) && check[1] {
		return !check[1], fmt.Errorf("cannot put your own king into check")
	}

	return true, nil
}

// Pawns, Knights, and Kings jump to a location (as opposed to crawling/sliding)
func validateMoveJump(m move) (bool, error) {
	// fmt.Println("coor", m.x1, " ", m.y1, " ", m.x2, " ", m.y2)
	j := [2]int{m.x1 - m.x2, m.y1 - m.y2}
	for _, i := range getDirections(m.startPiece) {
		// fmt.Println("match", i, " ", j)
		if i == j {
			return true, nil
		}
	}
	return false, fmt.Errorf("%s cannot move there", string(m.startPiece))
}

// Queen, Bishops, and Rooks crawl/slide across the board
func validateMoveCrawl(m move) (bool, error) {
	for _, i := range getDirections(m.startPiece) {
		x, y := m.x1+i[0], m.y1+i[1]
		for inBounds(x, y) {
			// fmt.Println(x, y, m.x1, m.y1, m.x2, m.y2, i)
			if x == m.x2 && y == m.y2 {
				// We made it to the endPiece
				return true, nil
			} else if m.brd[x][y] != '-' {
				// We ran into another piece
				break
			}
			x += i[0]
			y += i[1]
		}
	}
	return false, fmt.Errorf("%s cannot move there", string(m.startPiece))
}

// #######################################################################
// (Section 2) Check and Checkmate #######################################
// #######################################################################

func inCheck(b board) ([2]bool, error) {

	// Get location of both kings
	wk, bk, err := findKings(b)
	if err != nil {
		return [2]bool{false, false}, fmt.Errorf("%w", err)
	}

	// Create moves against the kings
	mbk := move{x2: bk[0], y2: bk[1], white: true, endPiece: 'k', brd: b}
	mwk := move{x2: wk[0], y2: wk[1], white: false, endPiece: 'K', brd: b}

	// To determine if a kings is in check,
	// attempt to validate moves of every pieces against the enemy king
	var whiteCheck, blackCheck bool
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if unicode.IsLetter(b[x][y]) {
				if unicode.IsUpper(b[x][y]) && !blackCheck {
					mbk.x1, mbk.y1, mbk.startPiece = x, y, b[x][y]
					blackCheck, _ = validateMove(mbk)
				} else if unicode.IsLower(b[x][y]) && !whiteCheck {
					mwk.x1, mwk.y1, mwk.startPiece = x, y, b[x][y]
					whiteCheck, _ = validateMove(mwk)
				}
			}
		}
	}

	return [2]bool{whiteCheck, blackCheck}, nil
}

func inCheckmate(b board, kingcolor bool) bool {

	var m move
	tmpB := b

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if unicode.IsLetter(b[x][y]) &&
				(kingcolor && unicode.IsUpper(b[x][y]) || !kingcolor && unicode.IsLower(b[x][y])) {

				m = move{x1: x, y1: y, white: kingcolor, startPiece: b[x][y], brd: tmpB}
				for z := 0; z < 8; z++ {
					for w := 0; w < 8; w++ {
						m.x2, m.y2, m.endPiece = z, w, b[z][w]
						validMove, _ := validateMove(m)
						if validMove {
							tmpB[m.x1][m.y1], tmpB[m.x2][m.y2] = '-', tmpB[m.x1][m.y1]
							inCheck, _ := inCheck(tmpB)
							if kingcolor && !inCheck[0] {
								return false
							} else if !kingcolor && !inCheck[1] {
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
func reportCheckAndCheckmate(b board) (bool, error) {

	check, err := inCheck(b)
	if err != nil {
		return false, err
	}
	if check[0] {
		if inCheckmate(b, White) {
			b.printBoard()
			fmt.Println("White is in checkmate!")
			fmt.Println("Black wins!")
			return true, nil
		}
		fmt.Println("White is in check!")
	} else if check[1] {
		if inCheckmate(b, !White) {
			b.printBoard()
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

func findKings(b board) ([2]int, [2]int, error) {

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
