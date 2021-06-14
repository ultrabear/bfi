package compiler

import (
	"runtime"
	"strings"
)

func countoccurances(indata string) bool {
	return strings.Contains(indata, "><") || strings.Contains(indata, "<>") || strings.Contains(indata, "+-") || strings.Contains(indata, "-+") || strings.Contains(indata, "[]") || strings.Contains(indata, "[->-[-]<]")
}

func stripoccurances(indata string) string {
	r := strings.NewReplacer("><", "", "<>", "", "+-", "", "-+", "", "[]", "", "[->-[-]<]", "[-]>[-]<")
	return r.Replace(indata)
}

// Filter a chunk to remove any non brainfuck instructions
func filterchunk(s string, r chan string) {

	var b strings.Builder

	b.Grow(len(s))

	for _, in := range s {
		switch in {
		case '+', '-', '>', '<', '[', ']', ',', '.':
			b.WriteByte(byte(in))
		}
	}

	r <- b.String()
}

func filterbfc(indata string) string {

	// The following code is confusing and long but it simply allocates the entire
	// string into smaller parts and then throws them at threads to be stripped seperateley
	// this gives a considerable speed improvement

	ncpu := runtime.NumCPU()

	splitwidth := (len(indata) + (ncpu - 1)) / ncpu

	// If the split size is so small many threads will probably only slow us down
	if splitwidth < 100 {
		rch := make(chan string)
		go filterchunk(indata, rch)
		return <-rch
	}

	returns := make([]chan string, ncpu)

	// Send jobs to threads
	for i := 0; i < len(returns); i++ {
		returns[i] = make(chan string)
		nextwidth := (splitwidth * i) + splitwidth
		if nextwidth > len(indata) {
			nextwidth = len(indata)
		}
		go filterchunk((indata)[splitwidth*i:nextwidth], returns[i])
	}

	// Get data back from threads
	joindata := make([]string, 0, len(returns))
	for _, item := range returns {
		joindata = append(joindata, <-item)
	}
	return strings.Join(joindata, "")
}

func CompressBFC(indata string) string {
	export := filterbfc(indata)
	for countoccurances(export) {
		export = stripoccurances(export)
	}
	return export
}
