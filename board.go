package main

import (
	"fmt"
	"strings"
)

type Board [8][8]rune

func newBoard(b *Board, pos string) error {

	// s := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"

	s := strings.Split(pos, "/")
	if len(s) != 8 {
		return fmt.Errorf("invalid board")
	}

	x := 0
	for _, row := range s {
		for y, piece := range row {

			if strings.Contains("012345678", string(piece)) {
				y += int(piece - '0')
				for j := 0; j < y; j++ {
					b[x][j] = '-'
				}
			} else if strings.Contains("prnbqkPRNBQK", string(piece)) {
				b[x][y] = piece
			} else {
				return fmt.Errorf("invalid board")
			}
		}
		x += 1
	}

	return nil
}

func printBoard(b *Board) {
	fmt.Println("   _A_B_C_D_E_F_G_H_")
	for i := 0; i < 8; i++ {
		fmt.Print(8-i, " |")
		for j := 0; j < 8; j++ {
			fmt.Print(" ", string(b[i][j]))
		}
		fmt.Print(" |")
		fmt.Println()
	}
	fmt.Println("  |_________________|")
}
