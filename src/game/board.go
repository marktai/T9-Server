package game

import (
	"log"
	"strings"
)

// This is the class that represents the all the squares and who owns them
// It is used in Game
type Board [9]Box

// Converts from board object to an array of ints
func (b *Board) Compress() [9]uint {
	var comprBoard [9]uint
	for i, _ := range b {
		comprBoard[i] = b[i].Compress()
	}
	return comprBoard
}

// Converts from an array of ints to a board object
func (b *Board) Decompress(compressed [9]uint) {
	for i, _ := range compressed {
		b[i].Decompress(compressed[i])
	}
}

// Prints the board with representation
func (b *Board) Print() {
	log.Print("Board: \n" + b.String(true))
}

// Returns a string representation of the board
func (b *Board) String(translate bool) string {
	out := ""
	for _, row := range b.StringArray(translate) {
		out += row + "\n"
	}
	return out
}

// Returns a string array representation with every line a new entry
func (b *Board) StringArray(translate bool) []string {
	outArray := make([]string, 0)
	var retArrays [3][5]string

	boxOfBoxes := b.Box()
	boxOfBoxes.CheckOwned()

	for row := 0; row < 3; row++ {

		if row != 0 {
			if boxOfBoxes.Owned == 0 {
				outArray = append(outArray, (""))
				outArray = append(outArray, (strings.Repeat("-", 37)))
				outArray = append(outArray, (""))
			} else {
				outArray = append(outArray, (""))
				outArray = append(outArray, (strings.Repeat("-", 11) + symbol(boxOfBoxes.Owned, translate) + strings.Repeat("-", 13) + symbol(boxOfBoxes.Owned, translate) + strings.Repeat("-", 13)))
				outArray = append(outArray, (""))
			}
		}

		for col := 0; col < 3; col++ {
			retArrays[col] = b[3*row+col].StringArray(translate)
		}
		for boxRow := 0; boxRow < 5; boxRow++ {
			outArray = append(outArray, (retArrays[0][boxRow] + "  |  " + retArrays[1][boxRow] + "  |  " + retArrays[2][boxRow]))
		}
	}
	return outArray
}

// Returns a box object that represents the board with every square representing one of the boxes
// Basically reduces the whole board to one box
func (b *Board) Box() *Box {
	var retBox Box
	for i := 0; i < 9; i++ {
		retBox.Squares[i] = b[i].Owned
	}
	retBox.CheckOwned()
	return &retBox
}
