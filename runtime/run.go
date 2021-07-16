package runtime

import (
	"fmt"
	"github.com/ultrabear/bfi/constants"
	"os"
	"unsafe"
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

func IndexUint(slice []uint, i int) *uint {
	return (*uint)(unsafe.Pointer(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&slice))) + uintptr(i)*unsafe.Sizeof(uint(0))))
}

func (bfc *Brainfuck) RunUnsafe(intfuck []uint) {

	// Basically cursed, has a ~2.6% speed advantage over BrainFuck.Run
	// Reads from the slices underlying memory location using unsafe pointers
	// This lets it bypass bounds checking, which is what gives a minor speed bump
	// As noted this is unsafe and on a corrupt intfuck slice could read uninit memory instead of bounds check and panic
	// Only use this when you know the input is not corrupt, for development it is worth it to swap to BrianFuck.Run

	if len(intfuck) == 0 {
		return
	}

	// Mainloop over brainfuck
	for i := 0; i < len(intfuck); i++ {
		next := *IndexUint(intfuck, i)
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
				fmt.Println(constants.RuntimeOverflow)
				os.Exit(1)
			}
		case constants.I_DecP: // Manually inlined bfc.DecP
			bfc.pointer--
			if bfc.pointer < 0 {
				fmt.Println(constants.RuntimeUnderflow)
				os.Exit(1)
			}
		case constants.I_Read:
			bfc.Read()
		case constants.I_Write:
			bfc.Write()
		case constants.I_LStart:
			i++
			if bfc.CurUnsafe() == 0 {
				i = int(*IndexUint(intfuck, i))
			}
		case constants.I_LEnd:
			i++
			if bfc.CurUnsafe() != 0 {
				i = int(*IndexUint(intfuck, i))
			}
		case constants.I_IncBy:
			i++
			next = *IndexUint(intfuck, i)
			bfc.IncByUnsafe(next)
		case constants.I_DecBy:
			i++
			next = *IndexUint(intfuck, i)
			bfc.DecByUnsafe(next)
		case constants.I_IncPBy: // Manually inlined bfc.IncPBy
			i++
			bfc.pointer += int(*IndexUint(intfuck, i))
			if bfc.pointer >= len(bfc.buffer) {
				fmt.Println(constants.RuntimeOverflow)
				os.Exit(1)
			}
		case constants.I_DecPBy: // Manually inlined bfc.DecPBy
			i++
			bfc.pointer -= int(*IndexUint(intfuck, i))
			if bfc.pointer < 0 {
				fmt.Println(constants.RuntimeUnderflow)
				os.Exit(1)
			}
		}
	}
}
