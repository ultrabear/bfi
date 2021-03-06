package compiler

import (
	"fmt"
	"os"

	con "github.com/ultrabear/bfi/constants"
)

type looper struct {
	precompiled [][2]int
	startloc    []int
	outmap      map[int]int
}

// Looping section of CompileLoops
func (L *looper) innerCompileLoops() {

	for _, val := range L.precompiled {

		// Value is a [
		if val[1] == 7 {

			// Add index to lifo stack
			L.startloc = append(L.startloc, val[0])

		} else {

			// Exit if ][ has occurred
			if len(L.startloc) == 0 {
				fmt.Fprintln(os.Stderr, con.SyntaxEndBeforeStart)
				os.Exit(1)
			}

			// Add value
			LStop, LStart := val[0], L.startloc[len(L.startloc)-1]
			L.outmap[LStart] = LStop
			L.outmap[LStop] = LStart

			// Pop off lifo stack
			L.startloc = L.startloc[:len(L.startloc)-1]

		}

	}

}

// Compiles loops
func (L *looper) compileloops() map[int]int {

	// Init output slice
	L.outmap = make(map[int]int, len(L.precompiled))

	L.innerCompileLoops()

	return L.outmap

}

// embedJumpMap will embed the jump points into the intfuck stream
func embedJumpMap(intfuck []uint, jumpmap map[int]int, keepcompiled [][2]int) []uint {

	// Should extend to enough space for inline loop instructions
	// If not this will panic
	// If it does panic it means the allocator routine is malfunctioning
	// In the event of a panic here check compiler.ToIntFuck
	var totalstream []uint

	if l := len(intfuck) + len(keepcompiled); l <= cap(intfuck) {
		totalstream = intfuck[:l]
	} else {
		panic("bfi/compiler.GetJumpMap: Not enough space allocated for loop instructions")
	}

	// loop count
	lc := 0

	// Calculate indexes of all values and store them for later
	indexes := make(map[int]int, len(keepcompiled))
	for i, v := range keepcompiled {
		indexes[v[0]] = i
	}

	// Rshift and inplace place loop index jump points
	for i := 1; i <= len(intfuck); i++ {
		// Testing if this item is a loop item
		if len(keepcompiled) > 0 && keepcompiled[len(keepcompiled)-1][0] == len(intfuck)-i {
			// Pop item off once its used
			keepcompiled = keepcompiled[:len(keepcompiled)-1]

			// Get the jump point this will jump to
			v := jumpmap[len(intfuck)-i]

			// Add index value to account for extra space taken by previous loop counts in slice
			totalstream[len(totalstream)-i-lc] = uint(v + indexes[v] + 1)

			// We add to the loop counter every loop to not lose track of rshift indexing
			lc++
		}
		// Rshift next value
		totalstream[len(totalstream)-i-lc] = intfuck[len(intfuck)-i]
	}

	// The length of intfuck is changed, return the new slice header
	return totalstream
}

// GetJumpMap will embed jump points into the intfuck stream after calculating their indexes
func GetJumpMap(intfuck []uint, sizeof int) []uint {

	// Compile brainfuck loops (3 steps)
	loops := looper{ // 1. Create looper object to handle loops
		precompiled: make([][2]int, 0, sizeof),
		startloc:    make([]int, 0, (sizeof+1)/2),
	}

	for i := 0; i < len(intfuck); i++ { // 2. Add [ ] to list
		switch intfuck[i] {
		case con.InstrucIncBy, con.InstrucDecBy, con.InstrucIncPBy, con.InstrucDecPBy:
			// Skip over instructions with argument fields
			i++
		case con.InstrucLStart, con.InstrucLEnd:
			loops.precompiled = append(loops.precompiled, [2]int{i, int(intfuck[i])})
		}
	}

	jumpmap := loops.compileloops() // 3. Compile loops recursively

	// Embed jumpmap into stream and return stream header
	return embedJumpMap(intfuck, jumpmap, loops.precompiled)
}
