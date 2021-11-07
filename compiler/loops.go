package compiler

import (
	"fmt"
	con "github.com/ultrabear/bfi/constants"
	"os"
)

type Looper struct {
	precompiled [][2]int
	startloc    []int
	output      [][2]int
}

// Recursive section of CompileLoops
func (L *Looper) innerCompileLoops() {

	for len(L.precompiled) > 0 {

		// Pop(0) slice
		val := L.precompiled[0]
		L.precompiled = L.precompiled[1:]

		// Value is a [
		if val[1] == 7 {

			// Add index to lifo stack
			L.startloc = append(L.startloc, val[0])

			// Continue searching for [] pairs
			L.innerCompileLoops()

		} else {

			// Exit if ][ has occured
			if len(L.startloc) == 0 {
				fmt.Fprintln(os.Stderr, con.SyntaxEndBeforeStart)
				os.Exit(1)
			}

			// Add value and pop off lifo stack
			L.output = append(L.output, [2]int{val[0], L.startloc[len(L.startloc)-1]})
			L.startloc = L.startloc[:len(L.startloc)-1]

			return
		}

	}

}

// Recursively compiles loops
func (L *Looper) Compileloops() map[int]int {

	// Init output slice
	L.output = make([][2]int, 0, len(L.precompiled)/2)

	L.innerCompileLoops()

	datamap := make(map[int]int, len(L.output)*2)

	for _, v := range L.output {
		datamap[v[0]] = v[1]
		datamap[v[1]] = v[0]
	}

	return datamap

}

func GetJumpMap(intfuck []uint, sizeof int) []uint {

	// Compile brainfuck loops (3 steps)
	loops := Looper{ // 1. Create looper object to handle loops
		precompiled: make([][2]int, 0, sizeof),
		startloc:    make([]int, 0, (sizeof+1)/2),
	}

	for i := 0; i < len(intfuck); i++ { // 2. Add [ ] to list
		switch intfuck[i] {
		case con.I_IncBy, con.I_DecBy, con.I_IncPBy, con.I_DecPBy:
			// Skip over instructions with argument fields
			i++
		case con.I_LStart, con.I_LEnd:
			loops.precompiled = append(loops.precompiled, [2]int{i, int(intfuck[i])})
		}
	}

	// Store original data for later
	keepcompiled := loops.precompiled

	jumpmap := loops.Compileloops() // 3. Compile loops recursively

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
