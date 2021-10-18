// Rendering tool for intfuck to C
package render

import (
	"fmt"
	con "github.com/ultrabear/bfi/constants"
	"strings"
)

// #define ARRSIZE X
// char* arr
// long ptr

var cmapping = [...]string{
	con.I_Zero:   "arr[ptr] = 0;",
	con.I_Inc:    "arr[ptr]++;",
	con.I_Dec:    "arr[ptr]--;",
	con.I_IncP:   "ptr++;",
	con.I_DecP:   "ptr--;",
	con.I_Read:   "arr[ptr] = fgetc(stdin);",
	con.I_Write:  "fputc(arr[ptr], stdout);",
	con.I_LStart: "while (arr[ptr] != 0) {",
	con.I_LEnd:   "}",
	con.I_IncBy:  "arr[ptr] += %d;",
	con.I_DecBy:  "arr[ptr] -= %d;",
	con.I_IncPBy: "ptr += %d;",
	con.I_DecPBy: "ptr -= %d;",
}

type CIntFuck struct {
	Data []uint
	Len  int
}

func (CIF CIntFuck) String() string {

	cout := make([]string, 1, len(CIF.Data)+2)

	out := cout[:]

	for i := 0; i < len(CIF.Data); i++ {
		switch CIF.Data[i] {
		case con.I_Zero, con.I_Inc, con.I_Dec, con.I_IncP, con.I_DecP, con.I_Read, con.I_Write:
			out = append(out, cmapping[CIF.Data[i]])
		case con.I_LStart, con.I_LEnd:
			i++
			out = append(out, cmapping[CIF.Data[i-1]])
		case con.I_IncBy, con.I_DecBy, con.I_IncPBy, con.I_DecPBy:
			i++
			out = append(out, fmt.Sprintf(cmapping[CIF.Data[i-1]], CIF.Data[i]))
		}
	}

	cout[0] = fmt.Sprintf("#include <stdio.h>\n#define ARRSIZE %d\nstatic char arr[ARRSIZE] = {0,};\nstatic long ptr = 0;\nint main() {", CIF.Len)

	cout = cout[:len(out)+1]

	cout[len(cout)-1] = "}"

	return strings.Join(cout, "\n")

}
