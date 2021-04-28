package runtime

func (bfc *Brainfuck) Run(intfuck []uint, jumpmap map[int]int) {

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
