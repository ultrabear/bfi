package compiler

import (
	"fmt"
	"github.com/ultrabear/bfi/constants"
	"os"
)

type Looper struct {
	precompiled [][2]int
	startloc    []int
}

func (L *Looper) Compileloops() map[int]int {
	datamap := map[int]int{}
	for len(L.precompiled) > 0 {
		if L.precompiled[0][1] == 7 {
			L.startloc = append(L.startloc, L.precompiled[0][0])
			L.precompiled = L.precompiled[1:]
			for k, v := range L.Compileloops() {
				datamap[k] = v
			}
		} else {
			if len(L.startloc) == 0 {
				fmt.Println(constants.SyntaxEndBeforeStart)
				os.Exit(1)
			}
			datamap[L.precompiled[0][0]] = L.startloc[len(L.startloc)-1]
			L.precompiled = L.precompiled[1:]
			L.startloc = L.startloc[:len(L.startloc)-1]
			return datamap
		}
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
		case 9, 10, 11, 12: // Skip over special instructions
			i++
		case 7, 8:
			loops.precompiled = append(loops.precompiled, [2]int{i, int(intfuck[i])})
		}
	}

	// Store original data for later
	keepcompiled := loops.precompiled

	jumpmap := loops.Compileloops() // 3. Compile loops recursively
	for k, v := range jumpmap {
		jumpmap[v] = k
	}

	// Should extend to enough space for inline loop instructions
	// If not this will panic
	// If it does panic it means the allocator routine is malfunctioning
	// In the event of a panic here check compiler.ToIntFuck
	totalstream := intfuck[:len(intfuck)+len(keepcompiled)]

	// loop count
	lc := 0

	// Calculate indexes of all values and store them for later
	indexes := map[int]int{}
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
			totalstream[len(totalstream)-i-lc] = uint(v+indexes[v]+1)

			// We add to the loop counter every loop to not lose track of rshift indexing
			lc++
		}
		// Rshift next value
		totalstream[len(totalstream)-i-lc] = intfuck[len(intfuck)-i]
	}

	// The length of intfuck is changed, return the new slice header
	return totalstream
}
