package main

import (
	"bytes"
	"fmt"
	"github.com/ultrabear/bfi/compiler"
	"github.com/ultrabear/bfi/constants"
	"github.com/ultrabear/bfi/render"
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

var (
	LStartByte = [1]byte{'['}
	LEndByte   = [1]byte{']'}
)

func ustring(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}

func RunCompile(indata []byte) ([]byte, []uint) {

	// Check amount of loops
	if bytes.Count(indata, LStartByte[:]) != bytes.Count(indata, LEndByte[:]) {
		fmt.Println(constants.SyntaxUnbalanced)
		os.Exit(1)
	}

	// Compress brainfuck and run static optimizations
	brainfuck := compiler.CompressBFC(indata)
	brainfuck = []byte(strings.NewReplacer("[-]", "0", "[+]", "0").Replace(ustring(brainfuck)))

	// Get count of loop items
	LoopCount := bytes.Count(brainfuck, LStartByte[:]) * 2

	// Convert to intfuck and optimize
	intfuck := compiler.PMoptimize(compiler.ToIntfuck(brainfuck, LoopCount))
	intfuck = compiler.GetJumpMap(intfuck, LoopCount)

	return brainfuck, intfuck
}

func RunFull(indata []byte) {

	brainfuck, intfuck := RunCompile(indata)

	// Instantize brainfuck execution environment
	bfc := runtime.Initbfc(max(bytes.Count(brainfuck, []byte{'>'})+1, 30000))

	// Run brainfuck
	bfc.RunUnsafe(intfuck)
}

func main() {

	var indata []byte

	if len(os.Args) > 2 && strings.Contains(os.Args[1], "f") {
		cont, readerr := os.ReadFile(os.Args[2])
		if readerr != nil {
			fmt.Println(constants.Error+"Could not open file:", os.Args[2])
			os.Exit(1)
		}
		indata = cont
	} else {
		indata = []byte(strings.Join(os.Args[1:], ""))
	}

	if len(os.Args) > 1 && strings.Contains(os.Args[1], "r") {
		_, intfuck := RunCompile(indata)
		fmt.Println(render.StrIntFuck(intfuck))
		os.Exit(0)
	}

	RunFull(indata)
}
