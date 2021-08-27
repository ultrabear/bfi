package main

import "testing"
import "os"

func TestNormal(T *testing.T) {

	files, e := os.ReadDir("./examples/")
	if e != nil {
		T.Error(e)
	}

	for _, file := range files {
		s, e := os.ReadFile("./examples/" + file.Name())
		if e != nil {
			T.Error(e)
		}
		RunFull(s)
	}

}
