package render

import (
	"fmt"
	"io"
	"strings"

	con "github.com/ultrabear/bfi/constants"
)

// #define ARRSIZE X
// char* arr
// long ptr

const endl = "\n"

var cmapping = [...]string{
	con.InstrucZero:   "arr[ptr] = 0;" + endl,
	con.InstrucInc:    "arr[ptr]++;" + endl,
	con.InstrucDec:    "arr[ptr]--;" + endl,
	con.InstrucIncP:   "ptr++;" + endl,
	con.InstrucDecP:   "ptr--;" + endl,
	con.InstrucRead:   "arr[ptr] = fgetc(stdin);" + endl,
	con.InstrucWrite:  "fputc(arr[ptr], stdout);" + endl,
	con.InstrucLStart: "while (arr[ptr] != 0) {" + endl,
	con.InstrucLEnd:   "}" + endl,
	con.InstrucIncBy:  "arr[ptr] += %d;" + endl,
	con.InstrucDecBy:  "arr[ptr] -= %d;" + endl,
	con.InstrucIncPBy: "ptr += %d;" + endl,
	con.InstrucDecPBy: "ptr -= %d;" + endl,
}

var bmapping = [len(cmapping)][]byte{}

func init() {
	for i, v := range cmapping {
		bmapping[i] = []byte(v)
	}
}

// CIntFuck is a wrapper for transpiling intfuck to C source code
type CIntFuck struct {
	Data []uint
	Len  int
}

// WriteTo implements the io.WriterTo interface
func (CIF *CIntFuck) WriteTo(w io.Writer) (int64, error) {

	var total int64

	header := fmt.Sprintf(
		`#include <stdio.h>

#define ARRSIZE %d

static char arr[ARRSIZE] = {0,};
static long ptr = 0;

int main() {
`, CIF.Len)

	n, err := w.Write([]byte(header))

	total += int64(n)
	if err != nil {
		goto Fail
	}

	for i := 0; i < len(CIF.Data); i++ {
		v := CIF.Data[i]
		switch v {
		case con.InstrucZero, con.InstrucInc, con.InstrucDec, con.InstrucIncP, con.InstrucDecP, con.InstrucRead, con.InstrucWrite:
			n, err = w.Write(bmapping[v])

			total += int64(n)
			if err != nil {
				goto Fail
			}
		case con.InstrucLStart, con.InstrucLEnd:
			i++
			n, err = w.Write(bmapping[v])

			total += int64(n)
			if err != nil {
				goto Fail
			}
		case con.InstrucIncBy, con.InstrucDecBy, con.InstrucIncPBy, con.InstrucDecPBy:
			i++
			n, err = w.Write([]byte(fmt.Sprintf(cmapping[v], CIF.Data[i])))

			total += int64(n)
			if err != nil {
				goto Fail
			}
		}
	}

	n, err = w.Write([]byte("}\n"))
	total += int64(n)

Fail:
	return total, err
}

func (CIF *CIntFuck) String() string {

	var b strings.Builder

	_, _ = CIF.WriteTo(&b)

	return b.String()
}
