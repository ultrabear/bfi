package compiler

import (
	"fmt"
	"github.com/ultrabear/bfi/constants"
	"os"
)

type Looper struct {
	precompiled [][]int
	startloc    []int
}

func (L *Looper) Compileloops() map[int]int {
	datamap := map[int]int{}
	for len(L.precompiled) > 0 {
		if L.precompiled[0][1] == 7 {
			L.startloc = append(L.startloc, L.precompiled[0][0])
			// Yes I know its slow as shit to take the first item out of a array
			// But so far its not been a bottleneck and I cba making it better
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

func GetJumpMap(intfuck []uint, sizeof int) map[int]int {

	// Compile brainfuck loops (3 steps)
	loops := Looper{ // 1. Create looper object to handle loops
		precompiled: make([][]int, 0, sizeof),
		startloc:    make([]int, 0, sizeof),
	}

	for i := 0; i < len(intfuck); i++ { // 2. Add [ ] to list
		switch intfuck[i] {
		case 9, 10, 11, 12: // Skip over special instructions
			i++
		case 7, 8:
			loops.precompiled = append(loops.precompiled, []int{i, int(intfuck[i])})
		}
	}

	jumpmap := loops.Compileloops() // 3. Compile loops recursively
	for k, v := range jumpmap {
		jumpmap[v] = k
	}

	return jumpmap
}
