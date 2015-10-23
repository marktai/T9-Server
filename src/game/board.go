package game

import (
	"log"
	"strings"
)

type Board [9]Box

func (b *Board) Compress() [9]uint {
	var comprBoard [9]uint
	for i, _ := range b {
		comprBoard[i] = b[i].Compress()
	}
	return comprBoard
}

func (b *Board) Decompress(compressed [9]uint) {
	for i, _ := range compressed {
		b[i].Decompress(compressed[i])
	}
}

func (b *Board) Print() {
	log.Print("Board: \n" + b.String(true))
}

func (b *Board) String(translate bool) string {
	out := ""
	for _, row := range b.StringArray(translate) {
		out += row + "\n"
	}
	return out
}

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

func (b *Board) Box() *Box {
	var retBox Box
	for i := 0; i < 9; i++ {
		retBox.Squares[i] = b[i].Owned
	}
	retBox.CheckOwned()
	return &retBox
}
