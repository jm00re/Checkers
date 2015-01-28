package main

import (
	"bufio"
	"fmt"
	//"github.com/davecheney/profile"
	//"math"
	"os"
	//"strconv"
	"strings"
)

type Board struct {
	player     bool
	blackDiscs uint32
	whiteDiscs uint32
	blackKings uint32
	whiteKings uint32
	occupied   uint32
}

// Bitmasks
var removeRight uint32
var removeLeft uint32
var removeFront uint32
var removeBack uint32

var removeLeftTwo uint32
var removeRightTwo uint32

func main() {
	// Bitmask declaration
	removeLeft = 0xefefefef
	removeRight = 0xf7f7f7f7
	removeFront = 0xfffffff
	removeBack = 0xfffffff0
	removeRightTwo = 0x77777777
	removeLeftTwo = 0xeeeeeeee

	player, board := ReadBoard()
	b := GenerateBoard(player, board)
	//fmt.Printf("%b\n", b.blackDiscs)
	//PrintBoard(b)
	//fmt.Println()
	//PrintBitBoard(b.blackDiscs & removeLeft)
	//fmt.Println()
	//BlackDiscMoves(b)
	//fmt.Println()
	//WhiteDiscMoves(b)
	BlackDiscCaptures(b)
	fmt.Println()
}

func BlackDiscMoves(b Board) (moves uint32) {
	// For Left move:
	// Remove  leftmost discs (Can't move left)
	// Shift remaining discs left
	// AND NOT against occupied location (Can't move into occupied spot)
	// Repeat for right move
	// OR two groups together
	moves = (((b.blackDiscs & removeLeft) << 4) &^ b.occupied) | (((b.blackDiscs & removeRight) << 5) &^ b.occupied)
	return
}

func BlackDiscCaptures(b Board) (moves uint32) {
	// The first half check if an opposing piece is diagonal,
	// then check if there is a black space open beyond that
	// Need to reduce the shift count 1 one because you're moving one less position in the second row
	moves = (((((b.blackDiscs & removeRightTwo) << 5) & b.whiteDiscs) << 4) &^ b.occupied) |
		(((((b.blackDiscs & removeLeftTwo) << 4) & b.whiteDiscs) << 3) &^ b.occupied)
	return
}

func WhiteDiscMoves(b Board) (moves uint32) {
	// Reverse shifts from BlackDiscMoves
	moves = (((b.whiteDiscs & removeRight) >> 4) &^ b.occupied) | (((b.whiteDiscs & removeLeft) >> 5) &^ b.occupied)
	PrintBitBoard(moves)
	return
}

func PrintBitBoard(b uint32) {
	// Shifts by i and checks if the value is 1. If it is print an x to represent that it's filled
	var shift uint8
	for i := 0; i < 32; i++ {
		shift = uint8(i)
		if i%4 == 0 {
			fmt.Println()
		}
		if (i >= 0 && i <= 3) || (i >= 8 && i <= 11) || (i >= 16 && i <= 19) || (i >= 24 && i <= 27) {
			fmt.Print("_")
			if (b>>shift)&1 != 0 {
				fmt.Print("x")
			} else {
				fmt.Print("_")
			}
		} else {
			if (b>>shift)&1 != 0 {
				fmt.Print("x")
			} else {
				fmt.Print("_")
			}
			fmt.Print("_")
		}
	}
}

func PrintBoard(b Board) {
	// Shifts by i and checks if the value is 1.
	// Prints the correct indictor based on the bitboard used
	if b.player {
		fmt.Print("b")
	} else {
		fmt.Print("w")
	}
	var shift uint8
	for i := 0; i < 32; i++ {
		shift = uint8(i)
		if i%4 == 0 {
			fmt.Println()
		}
		if (i >= 0 && i <= 3) || (i >= 8 && i <= 11) || (i >= 16 && i <= 19) || (i >= 24 && i <= 27) {
			fmt.Print("_")
			if (b.blackDiscs>>shift)&1 != 0 {
				fmt.Print("b")
			} else if (b.whiteDiscs>>shift)&1 != 0 {
				fmt.Print("w")
			} else if (b.blackKings>>shift)&1 != 0 {
				fmt.Print("B")
			} else if (b.whiteKings>>shift)&1 != 0 {
				fmt.Print("W")
			} else {
				fmt.Print("_")
			}
		} else {
			if (b.blackDiscs>>shift)&1 != 0 {
				fmt.Print("b")
			} else if (b.whiteDiscs>>shift)&1 != 0 {
				fmt.Print("w")
			} else if (b.blackKings>>shift)&1 != 0 {
				fmt.Print("B")
			} else if (b.whiteKings>>shift)&1 != 0 {
				fmt.Print("W")
			} else {
				fmt.Print("_")
			}
			fmt.Print("_")
		}
	}
}

func GenerateBoard(player bool, board [64]uint8) (b Board) {
	// Read STDIN into a parseable format
	var reducedBoard [32]uint8
	for i := 0; i < 32; i++ {
		if (i >= 0 && i <= 3) || (i >= 8 && i <= 11) || (i >= 16 && i <= 19) || (i >= 24 && i <= 27) {
			reducedBoard[i] = board[(i*2)+1]
		} else {
			reducedBoard[i] = board[i*2]
		}
	}
	var shift uint8
	//fmt.Println(reducedBoard)
	for i := 0; i < 32; i++ {
		shift = uint8(i)
		if reducedBoard[i] == 98 {
			b.blackDiscs = b.blackDiscs ^ (1 << shift)
		} else if reducedBoard[i] == 119 {
			b.whiteDiscs = b.whiteDiscs ^ (1 << shift)
		} else if reducedBoard[i] == 66 {
			b.blackKings = b.blackKings ^ (1 << shift)
		} else if reducedBoard[i] == 87 {
			b.whiteKings = b.whiteKings ^ (1 << shift)
		}
	}
	b.player = player
	b.occupied = (b.blackDiscs | b.whiteDiscs | b.blackKings | b.whiteKings)
	return
}

func ReadBoard() (player bool, board [64]uint8) {
	scanner := bufio.NewScanner(os.Stdin)
	// Read player number
	scanner.Scan()
	tempPlayer := scanner.Text()
	tempPlayer = strings.TrimSpace(tempPlayer)
	//fmt.Println(tempPlayer)
	if tempPlayer == "b" {
		player = true
	} else {
		player = false
	}
	scanner.Scan()
	for r := 0; r < 8; r++ {
		scanner.Scan()
		for c := 0; c < 8; c++ {
			board[(r*8)+c] = scanner.Text()[c]
		}
	}
	return
}

// Utility Functions
func Max(a int32, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a int32, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}
