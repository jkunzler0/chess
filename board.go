package main

import (
	"fmt"
	"strings"
)

type Board [8][8]rune

func defaultBoard(b *Board) error {
	return newBoard(b, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
}

func newBoard(b *Board, pos string) error {

	// s := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"

	s := strings.Split(pos, "/")
	if len(s) != 8 {
		return fmt.Errorf("invalid board")
	}

	for y, rank := range s {
		x := 0
		for _, piece := range rank {
			if strings.Contains("12345678", string(piece)) {
				numOfSpaces := int(piece - '0')
				for j := x; j < x+numOfSpaces; j++ {
					b[j][y] = '-'
				}
				x += numOfSpaces - 1
			} else if strings.Contains("prnbqkPRNBQK", string(piece)) {
				b[x][y] = piece
			} else {
				return fmt.Errorf("board contains invalid char %s", string(piece))
			}
			x += 1
		}
		if x != 8 {
			return fmt.Errorf("board is missing pieces")
		}
	}
	return nil
}

func printBoardBasic(b *Board) {
	fmt.Println("   _A_B_C_D_E_F_G_H_")
	for i := 0; i < 8; i++ {
		fmt.Print(8-i, " |")
		for j := 0; j < 8; j++ {
			fmt.Print(" ", string(b[j][i]))
		}
		fmt.Println(" |")
	}
	fmt.Println("  |_________________|")
}

func printBoard(b *Board) {
	var characters = map[rune]string{'P': "\u2659", 'p': "\u265F",
		'N': "\u2658", 'n': "\u265E",
		'B': "\u2657", 'b': "\u265D",
		'R': "\u2656", 'r': "\u265C",
		'Q': "\u2655", 'q': "\u265B",
		'K': "\u2654", 'k': "\u265A",
		'-': "-"}

	fmt.Println("   _A_B_C_D_E_F_G_H_")
	for i := 0; i < 8; i++ {
		fmt.Print(8-i, " |")
		for j := 0; j < 8; j++ {
			fmt.Print(" ", characters[b[j][i]])
		}
		fmt.Println(" |")
	}
	fmt.Println("  |_________________|")
}
