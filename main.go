package main

import (
	"bufio"
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

// fatal will print to stderr and exit
func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
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
		fatal(constants.SyntaxUnbalanced)
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

type conf struct {
	n string
	b bool
}

type confContainer map[string]*conf

func (C confContainer) newConf(n string) *conf {
	c := &conf{n: n}
	C[n] = c
	return c
}

func (C confContainer) parse() {

	if len(os.Args) > 1 {

		for k, v := range C {
			v.b = strings.Contains(os.Args[1], k)
		}

	}
}

func getbrainfuckstring(b bool) []byte {

	var indata []byte

	if b {
		if len(os.Args) > 2 {
			cont, readerr := os.ReadFile(os.Args[2])
			if readerr != nil {
				fatal(constants.Error+"Could not open file:", os.Args[2])
			}
			indata = cont
		} else {
			fatal(constants.Error + "Could not open file:")
		}
	} else {
		indata = []byte(strings.Join(os.Args[1:], ""))
	}

	return indata
}

func renderbrainfuck(indata []byte, renderc bool) {

	brainfuck, intfuck := RunCompile(indata)

	if renderc {
		cintf := render.CIntFuck{
			Data: intfuck,
			Len:  max(bytes.Count(brainfuck, []byte{'>'})+1, 30000),
		}

		w := bufio.NewWriter(os.Stdout)
		cintf.WriteTo(w)
		w.Flush()

		return
	}

	fmt.Println(render.StrIntFuck(intfuck))

}

func main() {

	conf := make(confContainer)

	f := conf.newConf("f")
	r := conf.newConf("r")
	c := conf.newConf("c")

	conf.parse()

	indata := getbrainfuckstring(f.b)

	if r.b {
		renderbrainfuck(indata, c.b)
		os.Exit(0)
	}

	RunFull(indata)
}
