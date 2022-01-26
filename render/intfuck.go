// Package render is a collection of rendering tools for intfuck
// streams to other formats.
package render

import (
	"fmt"
	"strconv"
	"strings"

	con "github.com/ultrabear/bfi/constants"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
	colBlue  = "\033[94m"
	colRed   = "\033[91m"
	colCyan  = "\033[96m"
	colNone  = "\033[0m"
)

type instruc struct {
	name string
	col  string
}

var instrucs = [...]instruc{
	con.InstrucZero:   {"ZERO", colGreen},
	con.InstrucInc:    {"INC", colGreen},
	con.InstrucDec:    {"DEC", colGreen},
	con.InstrucIncP:   {"INCP", colBlue},
	con.InstrucDecP:   {"DECP", colBlue},
	con.InstrucRead:   {"READ", colRed},
	con.InstrucWrite:  {"WRITE", colRed},
	con.InstrucLStart: {"LSTRT", colCyan},
	con.InstrucLEnd:   {"LEND", colCyan},
	con.InstrucIncBy:  {"INCB", colGreen},
	con.InstrucDecBy:  {"DECB", colGreen},
	con.InstrucIncPBy: {"INCPB", colBlue},
	con.InstrucDecPBy: {"DECPB", colBlue},
}

// StrIntFuck is a wrapper for representing intfuck as a human readable string
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
			case con.InstrucZero, con.InstrucInc, con.InstrucDec, con.InstrucIncP, con.InstrucDecP, con.InstrucRead, con.InstrucWrite:
				v := instrucs[SIF[i]]
				b.WriteString(w.writeinst(v.name, v.col))

			case con.InstrucLStart, con.InstrucLEnd, con.InstrucIncBy, con.InstrucDecBy, con.InstrucIncPBy, con.InstrucDecPBy:
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
