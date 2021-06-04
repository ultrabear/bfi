package runtime

import (
	"unsafe"
)

func (bfc *Brainfuck) RunOld(intfuck []uint, jumpmap map[int]int) {

	// Slower than switch statements somehow, ill leave it cause its neat but what the hell

	funcmap := []func(){
		bfc.Zero,
		bfc.Inc,
		bfc.Dec,
		bfc.IncP,
		bfc.DecP,
		bfc.Read,
		bfc.Write,
	}

	optifuncmap := []func(uint){
		bfc.IncBy,
		bfc.DecBy,
		bfc.IncPBy,
		bfc.DecPBy,
	}

	// Mainloop over brainfuck
	for i := 0; i < len(intfuck); i++ {
		switch intfuck[i] {
		case 7:
			if bfc.Cur() == 0 {
				i = jumpmap[i]
			}
		case 8:
			if bfc.Cur() != 0 {
				i = jumpmap[i]
			}
		default:
			if intfuck[i] <= 6 {
				funcmap[intfuck[i]]()
			} else {
				optifuncmap[intfuck[i]-9](intfuck[i+1])
				i++
			}
		}
	}
}

func (bfc *Brainfuck) Run(intfuck []uint, jumpmap map[int]int) {

	// Can golang make it faster with pure switches?
	// Benchmarks say yes somehow, this language has some
	// stupid level optimizations that I will never understand

	// Mainloop over brainfuck
	for i := 0; i < len(intfuck); i++ {
		switch intfuck[i] {
		case 0:
			bfc.Zero()
		case 1:
			bfc.Inc()
		case 2:
			bfc.Dec()
		case 3:
			bfc.IncP()
		case 4:
			bfc.DecP()
		case 5:
			bfc.Read()
		case 6:
			bfc.Write()
		case 7:
			if bfc.Cur() == 0 {
				i = jumpmap[i]
			}
		case 8:
			if bfc.Cur() != 0 {
				i = jumpmap[i]
			}
		case 9:
			i++
			bfc.IncBy(intfuck[i])
		case 10:
			i++
			bfc.DecBy(intfuck[i])
		case 11:
			i++
			bfc.IncPBy(intfuck[i])
		case 12:
			i++
			bfc.DecPBy(intfuck[i])
		}
	}
}

func (bfc *Brainfuck) RunUnsafe(intfuck []uint, jumpmap map[int]int) {

	// Basically cursed, has a ~2.6% speed advantage over BrainFuck.Run
	// Reads from the slices underlying memory location using unsafe pointers
	// This lets it bypass bounds checking, which is what gives a minor speed bump
	// As noted this is unsafe and on a corrupt intfuck slice could read uninit memory instead of bounds check and panic
	// Only use this when you know the input is not corrupt, for development it is worth it to swap to BrianFuck.Run

	jumpBy := unsafe.Sizeof(uint(0))
	SlicePtr := unsafe.Pointer(&(intfuck[0]))

	// Mainloop over brainfuck
	for i := 0; i < len(intfuck); i++ {
		next := *(*uint)(unsafe.Pointer(uintptr(SlicePtr)+(uintptr(i)*jumpBy)))
		switch next {
		case 0:
			bfc.Zero()
		case 1:
			bfc.Inc()
		case 2:
			bfc.Dec()
		case 3:
			bfc.IncP()
		case 4:
			bfc.DecP()
		case 5:
			bfc.Read()
		case 6:
			bfc.Write()
		case 7:
			if bfc.Cur() == 0 {
				i = jumpmap[i]
			}
		case 8:
			if bfc.Cur() != 0 {
				i = jumpmap[i]
			}
		case 9:
			i++
			next = *(*uint)(unsafe.Pointer(uintptr(SlicePtr)+(uintptr(i)*jumpBy)))
			bfc.IncBy(next)
		case 10:
			i++
			next = *(*uint)(unsafe.Pointer(uintptr(SlicePtr)+(uintptr(i)*jumpBy)))
			bfc.DecBy(next)
		case 11:
			i++
			next = *(*uint)(unsafe.Pointer(uintptr(SlicePtr)+(uintptr(i)*jumpBy)))
			bfc.IncPBy(next)
		case 12:
			i++
			next = *(*uint)(unsafe.Pointer(uintptr(SlicePtr)+(uintptr(i)*jumpBy)))
			bfc.DecPBy(next)
		}
	}
}
