package game

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

// Turn
// 0_ -> player 1 turn
// _0-8 -> box 0-8, 9 -> anywhere

type Box struct {
	Owned   uint
	Squares [9]uint
}

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

func (b *Box) Compress() uint {

	var total uint

	total += b.Owned
	for i := 8; i >= 0; i-- {
		total *= 4
		total += b.Squares[i]
	}

	return total
}

func (b *Box) Print() {

	out := fmt.Sprintln("Box:")
	out += b.String(true)
	log.Println(out)
}

func (b *Box) String(translate bool) string {
	out := ""
	retArray := b.StringArray(translate)

	for i := 0; i < 5; i++ {
		out += retArray[i] + "\n"
	}
	return out
}

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

func tripEqualityAndNot0(a, b, c uint) bool {
	return a != 0 && a == b && a == c
}

func (b *Box) CheckOwned() {

	//horizontal
	for i := 0; i < 3; i++ {
		if tripEqualityAndNot0(b.Squares[3*i], b.Squares[3*i+1], b.Squares[3*i+2]) {
			b.Owned = b.Squares[3*i]
			return
		}
	}

	//vertical
	for i := 0; i < 3; i++ {
		if tripEqualityAndNot0(b.Squares[i], b.Squares[i+3], b.Squares[i+6]) {
			b.Owned = b.Squares[i]
			return
		}
	}

	if tripEqualityAndNot0(b.Squares[0], b.Squares[4], b.Squares[8]) {
		b.Owned = b.Squares[0]
		return
	}

	if tripEqualityAndNot0(b.Squares[2], b.Squares[4], b.Squares[6]) {
		b.Owned = b.Squares[2]
		return
	}
}
