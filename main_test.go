package main

import "testing"
import "os"

func TestRun(t *testing.T) {

	f, e := os.ReadFile("./examples/beemovie.bf")
	if e != nil {
		panic(e)
	}

	RunFull(string(f))

}
