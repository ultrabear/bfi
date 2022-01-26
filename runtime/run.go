package runtime

import (
	"fmt"
	"os"

	"github.com/ultrabear/bfi/constants"
)

// Run will run brainfuck code
func (bfc *Brainfuck) Run(intfuck []uint) {

	// Can golang make it faster with pure switches?
	// Benchmarks say yes somehow, this language has some
	// stupid level optimizations that I will never understand

	// Update on ^ I know why its faster !!
	// Function inlining

	// Mainloop over brainfuck
	for i := 0; i < len(intfuck); i++ {
		switch intfuck[i] {
		case constants.InstrucZero:
			bfc.Zero()
		case constants.InstrucInc:
			bfc.Inc()
		case constants.InstrucDec:
			bfc.Dec()
		case constants.InstrucIncP:
			bfc.IncP()
		case constants.InstrucDecP:
			bfc.DecP()
		case constants.InstrucRead:
			bfc.Read()
		case constants.InstrucWrite:
			bfc.Write()
		case constants.InstrucLStart:
			i++
			if bfc.Cur() == 0 {
				i = int(intfuck[i])
			}
		case constants.InstrucLEnd:
			i++
			if bfc.Cur() != 0 {
				i = int(intfuck[i])
			}
		case constants.InstrucIncBy:
			i++
			bfc.IncBy(intfuck[i])
		case constants.InstrucDecBy:
			i++
			bfc.DecBy(intfuck[i])
		case constants.InstrucIncPBy:
			i++
			bfc.IncPBy(intfuck[i])
		case constants.InstrucDecPBy:
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
		case constants.InstrucZero:
			bfc.ZeroUnsafe()
		case constants.InstrucInc:
			bfc.IncUnsafe()
		case constants.InstrucDec:
			bfc.DecUnsafe()
		case constants.InstrucIncP: // Manually inlined bfc.IncP
			bfc.pointer++
			if bfc.pointer >= len(bfc.buffer) {
				pquit(constants.RuntimeOverflow)
			}
		case constants.InstrucDecP: // Manually inlined bfc.DecP
			bfc.pointer--
			if bfc.pointer < 0 {
				pquit(constants.RuntimeUnderflow)
			}
		case constants.InstrucRead:
			bfc.Read()
		case constants.InstrucWrite:
			bfc.Write()
		case constants.InstrucLStart:
			i++
			if bfc.CurUnsafe() == 0 {
				i = int(intfuck[i])
			}
		case constants.InstrucLEnd:
			i++
			if bfc.CurUnsafe() != 0 {
				i = int(intfuck[i])
			}
		case constants.InstrucIncBy:
			i++
			next = intfuck[i]
			bfc.IncByUnsafe(next)
		case constants.InstrucDecBy:
			i++
			next = intfuck[i]
			bfc.DecByUnsafe(next)
		case constants.InstrucIncPBy: // Manually inlined bfc.IncPBy
			i++
			bfc.pointer += int(intfuck[i])
			if bfc.pointer >= len(bfc.buffer) {
				pquit(constants.RuntimeOverflow)
			}
		case constants.InstrucDecPBy: // Manually inlined bfc.DecPBy
			i++
			bfc.pointer -= int(intfuck[i])
			if bfc.pointer < 0 {
				pquit(constants.RuntimeUnderflow)
			}
		}
	}
}
