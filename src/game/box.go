package game

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

// This is the class that is used to represent 9 squares
// It is used in the Board class

// Turn
// 0_ -> player 1 turn
// _0-8 -> box 0-8, 9 -> anywhere

type Box struct {
	Owned   uint
	Squares [9]uint
}

// Sets the current Box to be equal to the decompressed value of input
func (b *Box) Decompress(compressed uint) error {
	for i := 0; i <= 9; i++ {
		num := compressed & 0x03
		if i == 9 {
			b.Owned = num
		} else {
			b.Squares[i] = num
		}
		compressed >>= 2
	}
	return nil
}

// Returns an int that represents the Box
func (b *Box) Compress() uint {
	var total uint

	total += b.Owned
	for i := 8; i >= 0; i-- {
		total *= 4
		total += b.Squares[i]
	}

	return total
}

// Prints out a spring representation of the Box
func (b *Box) Print() {
	out := fmt.Sprintln("Box:")
	out += b.String(true)
	log.Println(out)
}

// Returns a string representation of the Box
func (b *Box) String(translate bool) string {
	out := ""
	retArray := b.StringArray(translate)

	for i := 0; i < 5; i++ {
		out += retArray[i] + "\n"
	}
	return out
}

// Returns a string representation with every line a new entry
func (b *Box) StringArray(translate bool) [5]string {
	var retArray [5]string

	for i := 0; i < 3; i++ {
		if i != 0 {
			if b.Owned == 0 {
				retArray[2*i-1] = strings.Repeat("-", 9)
			} else {
				retArray[2*i-1] = fmt.Sprintf("--%s---%s--", symbol(b.Owned, translate), symbol(b.Owned, translate))
			}
		}
		retArray[2*i] = fmt.Sprintf("%s | %s | %s", symbol(b.Squares[3*i], translate), symbol(b.Squares[3*i+1], translate), symbol(b.Squares[3*i+2], translate))
	}
	return retArray
}

// Makes a move by player at square
func (b *Box) MakeMove(player, square uint) error {
	if player != 1 && player != 2 {
		return errors.New(fmt.Sprintf("%d is an invalid player", player))
	}
	if b.Squares[square] != 0 {
		return errors.New("Square already taken")
	}
	b.Squares[square] = player
	b.CheckOwned()
	return nil
}

// Check whether all 3 arguments are equal and not 0
// Used to solve whether a Box is owned
func tripEqualityAndNot0(a, b, c uint) bool {
	return a != 0 && a == b && a == c
}

// Checks whether this Box is owned
// Inherently prefers some configurations over others if there are two winners
func (b *Box) CheckOwned() uint {
	//horizontal
	for i := 0; i < 3; i++ {
		if tripEqualityAndNot0(b.Squares[3*i], b.Squares[3*i+1], b.Squares[3*i+2]) {
			b.Owned = b.Squares[3*i]
			return b.Owned
		}
	}

	//vertical
	for i := 0; i < 3; i++ {
		if tripEqualityAndNot0(b.Squares[i], b.Squares[i+3], b.Squares[i+6]) {
			b.Owned = b.Squares[i]
			return b.Owned
		}
	}

	if tripEqualityAndNot0(b.Squares[0], b.Squares[4], b.Squares[8]) {
		b.Owned = b.Squares[0]
		return b.Owned
	}

	if tripEqualityAndNot0(b.Squares[2], b.Squares[4], b.Squares[6]) {
		b.Owned = b.Squares[2]
		return b.Owned
	}

	return b.Owned
}
