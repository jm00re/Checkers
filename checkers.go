package main

import (
	"bufio"
	"fmt"
	//"github.com/davecheney/profile"
	//"math"
	"os"
	"strconv"
	"strings"
)

// TODO To calc the actual move xor new board with old board.

// I think I'm going to add a lot more bitboards to this:
// Even and Odd row variants of all the pieces
// All of the black/white pieces
// The pieces without the far left or right black squares
// Calculating each bitboard once per board and just looking them up might be way more efficient
type Board struct {
	player bool

	blackDiscs  uint32
	blackKings  uint32
	blackPieces uint32

	whiteDiscs  uint32
	whiteKings  uint32
	whitePieces uint32

	occupied uint32
}

// Bitmasks
const removeLeft uint32 = 0xefefefef
const removeRight uint32 = 0xf7f7f7f7
const removeRightTwo uint32 = 0x77777777
const removeLeftTwo uint32 = 0xeeeeeeee
const removeBack uint32 = 0xfffffff
const removeFront uint32 = 0xfffffff0
const keepFront uint32 = 0xf
const keepBack uint32 = 0xf0000000
const evenRows uint32 = 0xf0f0f0f
const oddRows uint32 = 0xf0f0f0f0

// Bitscan msv lookup table
var bitscanLookup = [32]uint8{0, 9, 1, 10, 13, 21, 2, 29, 11, 14, 16, 18, 22, 25, 3, 30, 8, 12, 20, 28, 15, 17, 24, 7, 19, 27, 23, 6, 26, 5, 4, 31}

func main() {
	// Bitmask declaration
	player, board := ReadBoard()
	b := GenerateBoard(player, board)
	PrintOn(b.blackDiscs)
	PrintBoardWithBitBoard(b, BlackDiscMoves(b))
	fmt.Println()
	fmt.Println(strconv.FormatInt(int64(b.blackDiscs), 2))
	fmt.Println("DownRight")
	PrintBitBoard(DownRightMoveSource(b.blackDiscs, BlackDiscMoves(b)))
	fmt.Println("DownLeft")
	PrintBitBoard(DownLeftMoveSource(b.blackDiscs, BlackDiscMoves(b)))
	fmt.Println()
	fmt.Println()
	PrintBoardWithBitBoard(b, WhiteDiscMoves(b))
	fmt.Println()
	fmt.Println(strconv.FormatInt(int64(b.whiteDiscs), 2))
	fmt.Println("UpRight")
	PrintBitBoard(UpRightMoveSource(b.whiteDiscs, WhiteDiscMoves(b)))
	fmt.Println("UpLeft")
	PrintBitBoard(UpLeftMoveSource(b.whiteDiscs, WhiteDiscMoves(b)))
	fmt.Println()
	temp := UpLeftMoveSource(b.whiteDiscs, WhiteDiscMoves(b))
	PrintBitBoard(MoveUpLeft(Bitscan(temp), temp))
	fmt.Println()
	temp2 := UpRightMoveSource(b.whiteDiscs, WhiteDiscMoves(b))
	PrintBitBoard(MoveUpRight(Bitscan(temp2), temp2))
	fmt.Println()
}

func MoveDownRight(move uint8, bb uint32) uint32 {
	return (((((1 << move) & oddRows) << 4) |
		(((1 << move) & evenRows & removeRight) << 5)) | bb) ^ (1 << move)
}

func MoveDownLeft(move uint8, bb uint32) uint32 {
	return (((((1 << move) & evenRows) << 4) |
		(((1 << move) & oddRows & removeLeft) << 3)) | bb) ^ (1 << move)
}

func MoveUpLeft(move uint8, bb uint32) uint32 {
	return (((((1 << move) & oddRows & removeLeft) >> 5) |
		(((1 << move) & evenRows) >> 4)) | bb) ^ (1 << move)

}
func MoveUpRight(move uint8, bb uint32) uint32 {
	return (((((1 << move) & evenRows & removeRight) >> 3) |
		(((1 << move) & oddRows) >> 4)) | bb) ^ (1 << move)
}

// There's probably a better way to do this. Just going to give this a try.
// Take two uint32 (the potential moves and were the pieces are)
// shift the move backwards and & with the board to produce a new bitboard with only black pieces
// that can make moves
// This new bitboard can be bitscanned to actually make the moves for the AI
// bb1 is board, bb2 is move bitboard
func DownRightMoveSource(bb1 uint32, bb2 uint32) uint32 {
	return (((bb2 & oddRows & removeLeft) >> 5) |
		((bb2 & evenRows) >> 4)) & bb1
}

func DownLeftMoveSource(bb1 uint32, bb2 uint32) uint32 {
	return (((bb2 & evenRows & removeRight) >> 3) |
		((bb2 & oddRows) >> 4)) & bb1
}

func UpLeftMoveSource(bb1 uint32, bb2 uint32) uint32 {
	return (((bb2 & oddRows) << 4) |
		((bb2 & evenRows & removeRight) << 5)) & bb1
}

func UpRightMoveSource(bb1 uint32, bb2 uint32) uint32 {
	return (((bb2 & evenRows) << 4) |
		((bb2 & oddRows & removeLeft) << 3)) & bb1
}

func BlackDiscMoves(b Board) uint32 {
	// So basically I'm taking advantage of masking the off rows,
	// Performing a shift and then doing it again for the other rows
	// After all the potential moves have been calc'd I mask the taken spaces
	// I messed this up orginally and didn't realize I needed to shift the rows differently
	return (((b.blackDiscs & removeLeft & oddRows) << 3) |
		((b.blackDiscs & oddRows) << 4) |
		((b.blackDiscs & evenRows) << 4) |
		((b.blackDiscs & removeRight & evenRows) << 5)) &^ b.occupied
}

func BlackDiscCaptures(b Board) uint32 {
	// The first half check if an opposing piece is diagonal,
	// then check if there is a black space open beyond that
	return (((((b.blackDiscs & removeLeftTwo & evenRows) << 4) &
		(b.whitePieces)) << 3) |
		((((b.blackDiscs & removeLeftTwo & oddRows) << 3) &
			(b.whitePieces)) << 4) |
		((((b.blackDiscs & removeRightTwo & evenRows) << 5) &
			(b.whitePieces)) << 4) |
		((((b.blackDiscs & removeRightTwo & oddRows) << 4) &
			(b.whitePieces)) << 5)) &^ b.occupied
}

func WhiteDiscMoves(b Board) uint32 {
	// Same dealie as BlackDiscMoves
	return (((b.whiteDiscs & oddRows) >> 4) |
		((b.whiteDiscs & removeLeft & oddRows) >> 5) |
		((b.whiteDiscs & removeRight & evenRows) >> 3) |
		((b.whiteDiscs & evenRows) >> 4)) &^ b.occupied
}

func WhiteDiscCaptures(b Board) uint32 {
	// Same dealie as BlackDiscCaptures
	return (((((b.whiteDiscs & removeLeftTwo & evenRows) >> 4) &
		(b.blackPieces)) >> 5) |
		((((b.whiteDiscs & removeLeftTwo & oddRows) >> 5) &
			(b.blackPieces)) >> 4) |
		((((b.whiteDiscs & removeRightTwo & evenRows) >> 3) &
			(b.blackPieces)) >> 4) |
		((((b.whiteDiscs & removeRightTwo & oddRows) >> 4) &
			(b.blackPieces)) >> 3)) &^ b.occupied
}

func BlackKingMoves(b Board) uint32 {
	if b.blackKings != 0 {
		return (((b.blackKings & evenRows) << 4) |
			((b.blackKings & removeRight & evenRows) << 5) |
			((b.blackKings & removeRight & evenRows) >> 3) |
			((b.blackKings & evenRows) >> 4) |
			((b.blackKings & removeLeft & oddRows) << 3) |
			((b.blackKings & oddRows) << 4) |
			((b.blackKings & oddRows) >> 4) |
			((b.blackKings & removeLeft & oddRows) >> 5)) &^ b.occupied
	}
	return 0
}

func WhiteKingMoves(b Board) uint32 {
	if b.whiteKings != 0 {
		return (((b.whiteKings & evenRows) << 4) |
			((b.whiteKings & removeRight & evenRows) << 5) |
			((b.whiteKings & removeRight & evenRows) >> 3) |
			((b.whiteKings & evenRows) >> 4) |
			((b.whiteKings & removeLeft & oddRows) << 3) |
			((b.whiteKings & oddRows) << 4) |
			((b.whiteKings & oddRows) >> 4) |
			((b.whiteKings & removeLeft & oddRows) >> 5)) &^ b.occupied
	}
	return 0
}

func NewBlackKings(b uint32) uint32 {
	return b & keepBack
}

func NewWhiteKings(b uint32) uint32 {
	return b & keepFront
}

// I straight up stole this from Kim Walisch ala https://chessprogramming.wikispaces.com/Kim+Walisch
// Returns the index of the lsb
func Bitscan(b uint32) uint8 {
	return bitscanLookup[((b^(b-1))*0x07C4ACDD)>>27]
}

func PrintOn(b uint32) {
	temp := Bitscan(b)
	for b != 0 {
		fmt.Println(temp)
		b ^= 1 << temp
		temp = Bitscan(b)
	}
}

// Returns the number of set bits in b
// Taken from wikipedia
func PopCount(b uint32) (count uint8) {
	for b != 0 {
		b &= (b - 1)
		count += 1
	}
	return
}

func PrintBitBoard(b uint32) {
	// Shifts by i and checks if the value is 1. If it is print an x to represent that it's filled
	var shift uint8
	for i := 0; i < 32; i++ {
		shift = uint8(i)
		if i%4 == 0 && i != 0 {
			fmt.Println()
		}
		if (i >= 0 && i <= 3) || (i >= 8 && i <= 11) ||
			(i >= 16 && i <= 19) || (i >= 24 && i <= 27) {
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
	fmt.Println()
}

func PrintBoardWithBitBoard(b Board, b2 uint32) {
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
		if i%4 == 0 && i != 0 {
			fmt.Println()
		}
		if (i >= 0 && i <= 3) || (i >= 8 && i <= 11) ||
			(i >= 16 && i <= 19) || (i >= 24 && i <= 27) {
			fmt.Print("_")
			if (b.blackDiscs>>shift)&1 != 0 {
				fmt.Print("b")
			} else if (b.whiteDiscs>>shift)&1 != 0 {
				fmt.Print("w")
			} else if (b.blackKings>>shift)&1 != 0 {
				fmt.Print("B")
			} else if (b.whiteKings>>shift)&1 != 0 {
				fmt.Print("W")
			} else if (b2>>shift)&1 != 0 {
				fmt.Print("x")
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
			} else if (b2>>shift)&1 != 0 {
				fmt.Print("x")
			} else {
				fmt.Print("_")
			}
			fmt.Print("_")
		}
	}
	fmt.Print()
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
		if i%4 == 0 && i != 0 {
			fmt.Println()
		}
		if (i >= 0 && i <= 3) || (i >= 8 && i <= 11) ||
			(i >= 16 && i <= 19) || (i >= 24 && i <= 27) {
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
	fmt.Print()
}

func GenerateBoard(player bool, board [64]uint8) (b Board) {
	// Read STDIN into a parseable format
	var reducedBoard [32]uint8
	for i := 0; i < 32; i++ {
		if (i >= 0 && i <= 3) || (i >= 8 && i <= 11) ||
			(i >= 16 && i <= 19) || (i >= 24 && i <= 27) {
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
	b.blackPieces = (b.blackDiscs | b.blackKings)
	b.whitePieces = (b.whiteDiscs | b.whiteKings)
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
