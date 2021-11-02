// Rendering tool for intfuck streams
package render

import (
	"fmt"
	con "github.com/ultrabear/bfi/constants"
	"strconv"
	"strings"
)

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

type bwriter struct {
	i          *int
	padding    int
	intpadding string
	getindex   func(int) uint
}

func (b *bwriter) writeinst(inst, col string) string {
	return fmt.Sprintf("%s%0"+b.intpadding+"d:"+strings.Repeat("=", b.padding-len(inst))+"%s"+colNone, col, *(b.i), inst)
}

func (b *bwriter) writeval(col string) string {
	return fmt.Sprintf("%s%0"+b.intpadding+"d:%0"+strconv.Itoa(b.padding)+"d"+colNone, col, *(b.i), b.getindex(*(b.i)))
}

const (
	colGreen = "\033[92m"
	colBlue = "\033[94m"
	colRed = "\033[91m"
	colCyan = "\033[96m"
	colNone = "\033[0m"
)

type instruc struct {
	name string
	col string
}

var instrucs = [...]instruc{
	con.I_Zero:   {"ZERO", colGreen},
	con.I_Inc:    {"INC", colGreen},
	con.I_Dec:    {"DEC", colGreen},
	con.I_IncP:   {"INCP", colBlue},
	con.I_DecP:   {"DECP", colBlue},
	con.I_Read:   {"READ", colRed},
	con.I_Write:  {"WRITE", colRed},
	con.I_LStart: {"LSTRT", colCyan},
	con.I_LEnd:   {"LEND", colCyan},
	con.I_IncBy:  {"INCB", colGreen},
	con.I_DecBy:  {"DECB", colGreen},
	con.I_IncPBy: {"INCPB", colBlue},
	con.I_DecPBy: {"DECPB", colBlue},
}


type StrIntFuck []uint

func (SIF StrIntFuck) String() string {

	var b strings.Builder

	// Amount of 0 padding for ints
	ipad := len(strconv.Itoa(len(SIF)))
	padding := max(ipad, 5)

	// padding for item, ipad for index, +1 for :, +9 for color control codes, + len-1 for seperator space, +2 for []
	tlen := ((padding+ipad+1+9)*len(SIF) + len(SIF) - 1 + 2)
	b.Grow(tlen)

	w := bwriter{
		padding:    padding,
		intpadding: strconv.Itoa(ipad),
	}

	w.getindex = func(i int) uint { return SIF[i] }

	b.WriteByte('[')

	{

		i := 0
		w.i = &i

	for i = 0; i < len(SIF); i++ {
		if i != 0 {
			b.WriteByte(' ')
		}
		switch SIF[i] {
		case con.I_Zero, con.I_Inc, con.I_Dec, con.I_IncP, con.I_DecP, con.I_Read, con.I_Write:
			v := instrucs[SIF[i]]
			b.WriteString(w.writeinst(v.name, v.col))

		case con.I_LStart, con.I_LEnd, con.I_IncBy, con.I_DecBy, con.I_IncPBy, con.I_DecPBy:
			v := instrucs[SIF[i]]

			b.WriteString(w.writeinst(v.name, v.col))

			i++
			b.WriteByte(' ')

			b.WriteString(w.writeval(v.col))
		}
	}

	}

	b.WriteByte(']')

	return b.String()
}
