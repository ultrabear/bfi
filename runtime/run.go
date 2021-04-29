package runtime

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
