package compiler

import (
	con "github.com/ultrabear/bfi/constants"
)

// ToIntfuck converts a stream of brainfuck text into intfuck instructions
func ToIntfuck(bfc []byte, extra int) []uint {

	// Convert the brainfuck to indexes of a list
	// This lets it avoid hashing and converting in the mainloop
	indexer := [256]uint{
		'0': con.InstrucZero,
		'+': con.InstrucInc,
		'-': con.InstrucDec,
		'>': con.InstrucIncP,
		'<': con.InstrucDecP,
		',': con.InstrucRead,
		'.': con.InstrucWrite,
		'[': con.InstrucLStart,
		']': con.InstrucLEnd,
	}

	// Convert brainfuck string to intfuck
	// Extra is added for adding loop data inplace
	ints := make([]uint, 0, len(bfc)+extra)
	for i := 0; i < len(bfc); i++ {
		ints = append(ints, indexer[bfc[i]])
	}

	return ints
}

// PMoptimize optimizes out repeat instructions into single instructions
func PMoptimize(input []uint) []uint {
	newlist := input[:0]
	types := map[uint]uint{
		con.InstrucInc:  con.InstrucIncBy,
		con.InstrucDec:  con.InstrucDecBy,
		con.InstrucIncP: con.InstrucIncPBy,
		con.InstrucDecP: con.InstrucDecPBy,
	}
	for i := 0; i < len(input); i++ {

		// Test if input is in optimizable range
		if _, ok := types[input[i]]; ok {

			// Amount of times instruction repeats
			ctr := uint(1)

			// Test if cur and next input are the same and count till not true
			for i+1 < len(input) && input[i] == input[i+1] {
				i++
				ctr++
			}

			// Check if instruction occurs more than once
			if ctr == 1 {
				newlist = append(newlist, input[i])
			} else {
				newlist = append(newlist, types[input[i]], ctr)
			}

		} else {
			// Non optimizable instructions
			newlist = append(newlist, input[i])
		}

	} // End for loop

	return newlist
}
