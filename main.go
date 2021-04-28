package main

import (
	"os"
	"io"
	"strings"
	"regexp"
	"fmt"
)

// Define errors up here to avoid looking for them in the code
const (
	Error string = "\033[91mERROR:\033[0m "
	SyntaxEndBeforeStart string = Error + "BF Syntax: Loop end defined before loop start"
	SyntaxUnbalanced string = Error + "BF Syntax: Unbalanced loop statements"
	RuntimeUnderflow string = Error + "BF Runtime: Underflowed pointer location"
	RuntimeOverflow string = Error + "BF Runtime: Overflowed pointer location"
)

type Brainfuck struct {
	buffer []byte
	pointer int
	stdin io.Reader
	stdout io.Writer
}

func (bfc *Brainfuck) Inc() {
	bfc.buffer[bfc.pointer]++
}

func (bfc *Brainfuck) Dec() {
	bfc.buffer[bfc.pointer]--
}

func (bfc *Brainfuck) IncBy(val uint) {
	bfc.buffer[bfc.pointer] += byte(val)
}

func (bfc *Brainfuck) DecBy(val uint) {
	bfc.buffer[bfc.pointer] -= byte(val)
}

func (bfc *Brainfuck) IncP() {
	bfc.pointer += 1
	if bfc.pointer >= len(bfc.buffer) {
		fmt.Println(RuntimeOverflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) DecP() {
	bfc.pointer -= 1
	if bfc.pointer < 0 {
		fmt.Println(RuntimeUnderflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) IncPBy(amt uint) {
	bfc.pointer += int(amt)
	if bfc.pointer >= len(bfc.buffer) {
		fmt.Println(RuntimeOverflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) DecPBy(amt uint) {
	bfc.pointer -= int(amt)
	if bfc.pointer < 0 {
		fmt.Println(RuntimeUnderflow)
		os.Exit(1)
	}
}


func (bfc *Brainfuck) Write() {
	bfc.stdout.Write([]byte{bfc.buffer[bfc.pointer]})
}

func (bfc *Brainfuck) Read() {
	var indata = make([]byte, 1)
	amt, err := bfc.stdin.Read(indata)
	if amt != 1 {
		indata[0] = 0
	}
	if err != nil {
		indata[0] = 0
	}
	bfc.buffer[bfc.pointer] = indata[0]
}

func (bfc *Brainfuck) Zero() {
	bfc.buffer[bfc.pointer] = 0
}

func (bfc *Brainfuck) Cur() int {
	return int(bfc.buffer[bfc.pointer])
}

func initbfc(size int) Brainfuck {
	return Brainfuck{
		buffer: make([]byte, size),
		pointer: 0,
		stdin: os.Stdin,
		stdout: os.Stdout,
	}
}

type Looper struct {
	precompiled [][]int
	startloc []int
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
				fmt.Println(SyntaxEndBeforeStart)
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


func countoccurances(indata *string) int {
	outdata := 0
	outdata += strings.Count(*indata, "><") + strings.Count(*indata, "<>")
	outdata += strings.Count(*indata, "+-") + strings.Count(*indata, "-+")
	outdata += strings.Count(*indata, "[]")
	outdata += strings.Count(*indata, "[->-[-]<]")
	return outdata
}

func stripoccurances(indata *string) string {
	outdata := strings.Replace(strings.Replace(*indata, "><", "", -1), "<>", "", -1)
	outdata = strings.Replace(strings.Replace(outdata, "+-", "", -1), "-+", "", -1)
	outdata = strings.Replace(outdata, "[]", "", -1)
	outdata = strings.Replace(outdata, "[->-[-]<]", "[-]>[-]<", -1)
	return outdata
}

func filterbfc(indata *string) string {
	gex, _ := regexp.Compile("[\\[\\]\\-\\+\\>\\<\\,\\.]+")
	data := gex.FindAllStringSubmatch(*indata, -1)
	joindata := make([]string, 0)
	for _, item := range data {
		joindata = append(joindata, strings.Join(item, ""))
	}
	return strings.Join(joindata, "")
}

func bfccompress(indata string) string {
	export := filterbfc(&indata)
	for countoccurances(&export) != 0 {
		export = stripoccurances(&export)
	}
	return export
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

func main () {

	var indata string

	if len(os.Args) > 2 && os.Args[1] == "f" {
		cont, readerr := os.ReadFile(os.Args[2])
		if readerr != nil {
			fmt.Println(Error + "Could not open file:", os.Args[2])
			os.Exit(1)
		}
		indata = string(cont)
	} else {
		indata = strings.Join(os.Args[1:], "")
	}


	// Check amount of loops
	if strings.Count(indata, "[") != strings.Count(indata, "]") {
		fmt.Println(SyntaxUnbalanced)
		os.Exit(1)
	}

	// Compress brainfuck and run static optomizations
	brainfuck := bfccompress(indata)
	brainfuck = strings.Replace(strings.Replace(brainfuck, "[-]", "0", -1), "[+]", "0", -1)

	// Convert the brainfuck to indexes of a list
	// This lets it avoid hashing and converting in the mainloop
	indexer := map[byte]uint{
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
	intfuck := func(bfc string, indexer map[byte]uint) []uint {
		ints := make([]uint, len(bfc))
		for i := 0; i < len(bfc); i++ {
			ints[i] = indexer[(bfc)[i]]
		}
		return ints
	}(brainfuck, indexer)

	// Optomize with plus and minus remapping
	intfuck = PMoptimize(intfuck)

	// Compile brainfuck loops (3 steps)
	loops := Looper { // 1. Create looper object to handle loops
		precompiled: make([][]int, 0, len(intfuck)),
		startloc: make([]int, 0, len(intfuck)),
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

	// Instantize brainfuck execution environment
	bfc := initbfc(len(brainfuck) + 1)

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
