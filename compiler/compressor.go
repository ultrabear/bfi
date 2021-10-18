package compiler

import (
	"bytes"
	"runtime"
	"strings"
	"unsafe"
)

func stob(s string) []byte {
	return []byte(s)
}

func ustring(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}

func countoccurances(indata []byte) bool {
	return bytes.Contains(indata, stob("><")) || bytes.Contains(indata, stob("<>")) || bytes.Contains(indata, stob("+-")) || bytes.Contains(indata, stob("-+")) || bytes.Contains(indata, stob("[]")) || bytes.Contains(indata, stob("[->-[-]<]"))
}

func stripoccurances(indata []byte) []byte {
	replacements := []string{
		"><", "",
		"<>", "",
		"+-", "",
		"-+", "",
		"[]", "",
		"[->-[-]<]", "[-]>[-]<",
	}

	r := strings.NewReplacer(replacements...)
	return []byte(r.Replace(ustring(indata)))
}

// Filter a chunk to remove any non brainfuck instructions
func filterchunk(s []byte, r chan []byte) {

	ns := s[:0]

	for _, in := range s {
		switch in {
		case '+', '-', '>', '<', '[', ']', ',', '.':
			ns = append(ns, in)
		}
	}

	r <- ns

}

func filterbfc(indata []byte) []byte {

	// The following code is confusing and long but it simply allocates the entire
	// string into smaller parts and then throws them at threads to be stripped seperateley
	// this gives a considerable speed improvement

	ncpu := runtime.NumCPU()

	splitwidth := (len(indata) + (ncpu - 1)) / ncpu

	// If the split size is so small many threads will probably only slow us down
	if splitwidth < 100 {
		rch := make(chan []byte)
		go filterchunk(indata, rch)
		return <-rch
	}

	returns := make([]chan []byte, ncpu)

	// Send jobs to threads
	for i := 0; i < len(returns); i++ {
		returns[i] = make(chan []byte)
		nextwidth := (splitwidth * i) + splitwidth
		if nextwidth > len(indata) {
			nextwidth = len(indata)
		}
		go filterchunk((indata)[splitwidth*i:nextwidth], returns[i])
	}

	// Get data back from threads
	joindata := make([][]byte, 0, len(returns))
	for _, item := range returns {
		joindata = append(joindata, <-item)
	}
	return bytes.Join(joindata, []byte{})
}

func CompressBFC(indata []byte) []byte {
	export := filterbfc(indata)
	for countoccurances(export) {
		export = stripoccurances(export)
	}
	return export
}
