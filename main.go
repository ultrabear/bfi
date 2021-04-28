package main

import (
	"fmt"
	"github.com/ultrabear/bfi/compiler"
	"github.com/ultrabear/bfi/constants"
	"github.com/ultrabear/bfi/runtime"
	"os"
	"strings"
)

func main() {

	var indata string

	if len(os.Args) > 2 && os.Args[1] == "f" {
		cont, readerr := os.ReadFile(os.Args[2])
		if readerr != nil {
			fmt.Println(constants.Error+"Could not open file:", os.Args[2])
			os.Exit(1)
		}
		indata = string(cont)
	} else {
		indata = strings.Join(os.Args[1:], "")
	}

	// Check amount of loops
	if strings.Count(indata, "[") != strings.Count(indata, "]") {
		fmt.Println(constants.SyntaxUnbalanced)
		os.Exit(1)
	}

	// Compress brainfuck and run static optomizations
	brainfuck := compiler.CompressBFC(indata)
	brainfuck = strings.Replace(strings.Replace(brainfuck, "[-]", "0", -1), "[+]", "0", -1)

	intfuck := compiler.PMoptimize(compiler.ToIntfuck(brainfuck))
	jumpmap := compiler.GetJumpMap(intfuck)

	// Instantize brainfuck execution environment
	bfc := runtime.Initbfc(len(brainfuck) + 1)

	// Run brainfuck
	bfc.Run(intfuck, jumpmap)

}
