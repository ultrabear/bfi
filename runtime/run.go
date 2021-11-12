package runtime

import (
	"fmt"
	"os"

	"github.com/ultrabear/bfi/constants"
)

func (bfc *Brainfuck) Run(intfuck []uint) {

	// Can golang make it faster with pure switches?
	// Benchmarks say yes somehow, this language has some
	// stupid level optimizations that I will never understand

	// Update on ^ I know why its faster !!
	// Function inlining

	// Mainloop over brainfuck
	for i := 0; i < len(intfuck); i++ {
		switch intfuck[i] {
		case constants.I_Zero:
			bfc.Zero()
		case constants.I_Inc:
			bfc.Inc()
		case constants.I_Dec:
			bfc.Dec()
		case constants.I_IncP:
			bfc.IncP()
		case constants.I_DecP:
			bfc.DecP()
		case constants.I_Read:
			bfc.Read()
		case constants.I_Write:
			bfc.Write()
		case constants.I_LStart:
			i++
			if bfc.Cur() == 0 {
				i = int(intfuck[i])
			}
		case constants.I_LEnd:
			i++
			if bfc.Cur() != 0 {
				i = int(intfuck[i])
			}
		case constants.I_IncBy:
			i++
			bfc.IncBy(intfuck[i])
		case constants.I_DecBy:
			i++
			bfc.DecBy(intfuck[i])
		case constants.I_IncPBy:
			i++
			bfc.IncPBy(intfuck[i])
		case constants.I_DecPBy:
			i++
			bfc.DecPBy(intfuck[i])
		}
	}
}

// non inlined function to lower size of brainfuck interpreter loop code
func pquit(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

// RunUnsafe will run brainfuck code with
// removed bounds checks on array accesses.
// It also has manually inlined functions and as such is
// the fastest way to interpret intfuck in bfi
func (bfc *Brainfuck) RunUnsafe(intfuck []uint) {

	// Basically cursed, has a ~2.6% speed advantage over BrainFuck.Run
	// Reads from the slices underlying memory location using unsafe pointers
	// This lets it bypass bounds checking, which is what gives a minor speed bump

	// Mainloop over brainfuck
	for i := 0; i < len(intfuck); i++ {
		next := intfuck[i]
		switch next {
		case constants.I_Zero:
			bfc.ZeroUnsafe()
		case constants.I_Inc:
			bfc.IncUnsafe()
		case constants.I_Dec:
			bfc.DecUnsafe()
		case constants.I_IncP: // Manually inlined bfc.IncP
			bfc.pointer++
			if bfc.pointer >= len(bfc.buffer) {
				pquit(constants.RuntimeOverflow)
			}
		case constants.I_DecP: // Manually inlined bfc.DecP
			bfc.pointer--
			if bfc.pointer < 0 {
				pquit(constants.RuntimeUnderflow)
			}
		case constants.I_Read:
			bfc.Read()
		case constants.I_Write:
			bfc.Write()
		case constants.I_LStart:
			i++
			if bfc.CurUnsafe() == 0 {
				i = int(intfuck[i])
			}
		case constants.I_LEnd:
			i++
			if bfc.CurUnsafe() != 0 {
				i = int(intfuck[i])
			}
		case constants.I_IncBy:
			i++
			next = intfuck[i]
			bfc.IncByUnsafe(next)
		case constants.I_DecBy:
			i++
			next = intfuck[i]
			bfc.DecByUnsafe(next)
		case constants.I_IncPBy: // Manually inlined bfc.IncPBy
			i++
			bfc.pointer += int(intfuck[i])
			if bfc.pointer >= len(bfc.buffer) {
				pquit(constants.RuntimeOverflow)
			}
		case constants.I_DecPBy: // Manually inlined bfc.DecPBy
			i++
			bfc.pointer -= int(intfuck[i])
			if bfc.pointer < 0 {
				pquit(constants.RuntimeUnderflow)
			}
		}
	}
}
