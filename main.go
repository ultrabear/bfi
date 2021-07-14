package main

import (
	"fmt"
	"github.com/ultrabear/bfi/compiler"
	"github.com/ultrabear/bfi/constants"
	"github.com/ultrabear/bfi/runtime"
	"os"
	"strings"
	"unsafe"
)

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func ustring(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}

func RunFull(indata string) {

	// Check amount of loops
	if strings.Count(indata, "[") != strings.Count(indata, "]") {
		fmt.Println(constants.SyntaxUnbalanced)
		os.Exit(1)
	}

	// Compress brainfuck and run static optimizations
	brainfuck := compiler.CompressBFC(indata)
	brainfuck = strings.NewReplacer("[-]", "0", "[+]", "0").Replace(brainfuck)

	// Get count of loop items
	LoopCount := strings.Count(brainfuck, "[") * 2

	// Convert to intfuck and optimize
	intfuck := compiler.PMoptimize(compiler.ToIntfuck(brainfuck, LoopCount))
	intfuck = compiler.GetJumpMap(intfuck, LoopCount)

	// Instantize brainfuck execution environment
	bfc := runtime.Initbfc(max(strings.Count(brainfuck, ">")+1, 30000))

	// Run brainfuck
	bfc.RunUnsafe(intfuck)
}

func main() {

	var indata string

	if len(os.Args) > 2 && os.Args[1] == "f" {
		cont, readerr := os.ReadFile(os.Args[2])
		if readerr != nil {
			fmt.Println(constants.Error+"Could not open file:", os.Args[2])
			os.Exit(1)
		}
		indata = ustring(cont)
	} else {
		indata = strings.Join(os.Args[1:], "")
	}

	RunFull(indata)
}
