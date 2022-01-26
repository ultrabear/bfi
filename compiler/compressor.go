// Package compiler contains functions for compiling brainfuck to intfuck and optimizing
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

func hasoccurances(indata []byte) bool {
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
func filterchunk(s []byte, r *[]byte, fin chan struct{}) {

	ns := s[:0]

	for _, in := range s {
		switch in {
		case '+', '-', '>', '<', '[', ']', ',', '.':
			ns = append(ns, in)
		}
	}

	(*r) = ns
	fin <- struct{}{}

}

func filterbfc(indata []byte) []byte {

	// The following code is confusing and long but it simply allocates the entire
	// string into smaller parts and then throws them at threads to be stripped seperateley
	// this gives a considerable speed improvement

	ncpu := runtime.NumCPU()

	splitwidth := (len(indata) + (ncpu - 1)) / ncpu

	// If the split size is so small many threads will probably only slow us down
	if splitwidth < 100 {
		var rch []byte
		done := make(chan struct{})
		go filterchunk(indata, &rch, done)
		<-done
		return rch
	}

	status := make(chan struct{})
	await := 0
	returns := make([][]byte, ncpu)

	// Send jobs to threads
	for i := 0; i < len(returns); i++ {
		nextwidth := (splitwidth * i) + splitwidth
		if nextwidth > len(indata) {
			nextwidth = len(indata)
		}
		go filterchunk((indata)[splitwidth*i:nextwidth], &returns[i], status)
		await++
	}

	// Wait for data back from threads
	for ; await != 0; await-- {
		<-status
	}

	return bytes.Join(returns, []byte{})
}

// CompressBFC removes any non brainfuck characters and then removes opposing brainfuck instructions
func CompressBFC(indata []byte) []byte {
	export := filterbfc(indata)
	for hasoccurances(export) {
		export = stripoccurances(export)
	}
	return export
}
