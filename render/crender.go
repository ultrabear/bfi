// Rendering tool for intfuck to C
package render

import (
	"fmt"
	con "github.com/ultrabear/bfi/constants"
	"io"
	"strings"
)

// #define ARRSIZE X
// char* arr
// long ptr

const endl = "\n"

var cmapping = [...]string{
	con.I_Zero:   "arr[ptr] = 0;" + endl,
	con.I_Inc:    "arr[ptr]++;" + endl,
	con.I_Dec:    "arr[ptr]--;" + endl,
	con.I_IncP:   "ptr++;" + endl,
	con.I_DecP:   "ptr--;" + endl,
	con.I_Read:   "arr[ptr] = fgetc(stdin);" + endl,
	con.I_Write:  "fputc(arr[ptr], stdout);" + endl,
	con.I_LStart: "while (arr[ptr] != 0) {" + endl,
	con.I_LEnd:   "}" + endl,
	con.I_IncBy:  "arr[ptr] += %d;" + endl,
	con.I_DecBy:  "arr[ptr] -= %d;" + endl,
	con.I_IncPBy: "ptr += %d;" + endl,
	con.I_DecPBy: "ptr -= %d;" + endl,
}

var bmapping = make([][]byte, len(cmapping))

func init() {
	for i, v := range cmapping {
		bmapping[i] = []byte(v)
	}
}

type CIntFuck struct {
	Data []uint
	Len  int
}

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
		case con.I_Zero, con.I_Inc, con.I_Dec, con.I_IncP, con.I_DecP, con.I_Read, con.I_Write:
			n, err = w.Write(bmapping[v])

			total += int64(n)
			if err != nil {
				goto Fail
			}
		case con.I_LStart, con.I_LEnd:
			i++
			n, err = w.Write(bmapping[v])

			total += int64(n)
			if err != nil {
				goto Fail
			}
		case con.I_IncBy, con.I_DecBy, con.I_IncPBy, con.I_DecPBy:
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

	CIF.WriteTo(&b)

	return b.String()
}
