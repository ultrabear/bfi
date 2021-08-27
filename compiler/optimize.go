package compiler

func ToIntfuck(bfc []byte, extra int) []uint {

	// Convert the brainfuck to indexes of a list
	// This lets it avoid hashing and converting in the mainloop
	indexer := [256]uint{
		'0': 0,
		'+': 1,
		'-': 2,
		'>': 3,
		'<': 4,
		',': 5,
		'.': 6,
		'[': 7,
		']': 8,
	}

	// Convert brainfuck string to intfuck
	// Extra is added for adding loop data inplace
	ints := make([]uint, 0, len(bfc)+extra)
	for i := 0; i < len(bfc); i++ {
		ints = append(ints, indexer[bfc[i]])
	}

	return ints
}

func PMoptimize(input []uint) []uint {
	newlist := input[:0]
	types := map[uint]uint{
		1: 9,
		2: 10,
		3: 11,
		4: 12,
	}
	for i := 0; i < len(input); i++ {
		if _, ok := types[input[i]]; ok {
			ctr := uint(1)
			for ctr > 0 {
				if i+1 < len(input) && input[i] == input[i+1] {
					i++
					ctr++
				} else {
					if ctr == 1 {
						newlist = append(newlist, input[i])
					} else {
						newlist = append(newlist, types[input[i]], ctr)
					}
					ctr = 0
				}
			}
		} else {
			newlist = append(newlist, input[i])
		}
	}
	return newlist
}
