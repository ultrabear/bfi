// Rendering tool for intfuck streams
package render

import (
	"fmt"
	"github.com/ultrabear/bfi/constants"
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
	i          int
	padding    int
	intpadding string
	getindex   func(int) uint
}

func (b *bwriter) writeinst(inst, col string) string {
	return fmt.Sprintf("%s%0"+b.intpadding+"d:"+strings.Repeat("=", b.padding-len(inst))+"%s%s", b.pcol(col), b.i, inst, b.pcol(""))
}

func (b *bwriter) writeval(col string) string {
	return fmt.Sprintf("%s%0"+b.intpadding+"d:"+"%0"+strconv.Itoa(b.padding)+"d%s", b.pcol(col), b.i, b.getindex(b.i), b.pcol(""))
}

func (b *bwriter) pcol(col string) string {
	switch col {
	case "red":
		return "\033[91m"
	case "green":
		return "\033[92m"
	case "yellow":
		return "\033[93m"
	case "blue":
		return "\033[94m"
	case "magenta":
		return "\033[95m"
	case "cyan":
		return "\033[96m"
	case "white":
		return "\033[97m"
	default:
		return "\033[0m"
	}
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

	for i := 0; i < len(SIF); i++ {
		if i != 0 {
			b.WriteByte(' ')
		}
		w.i = i
		switch SIF[i] {
		case constants.I_Zero:
			b.WriteString(w.writeinst("ZERO", "green"))
		case constants.I_Inc:
			b.WriteString(w.writeinst("INC", "green"))
		case constants.I_Dec:
			b.WriteString(w.writeinst("DEC", "green"))
		case constants.I_IncP:
			b.WriteString(w.writeinst("INCP", "blue"))
		case constants.I_DecP:
			b.WriteString(w.writeinst("DECP", "blue"))
		case constants.I_Read:
			b.WriteString(w.writeinst("READ", "red"))
		case constants.I_Write:
			b.WriteString(w.writeinst("WRITE", "red"))
		case constants.I_LStart:
			b.WriteString(w.writeinst("LSTRT", "cyan"))
			i++
			w.i = i
			b.WriteByte(' ')
			b.WriteString(w.writeval("cyan"))
		case constants.I_LEnd:
			b.WriteString(w.writeinst("LEND", "cyan"))
			i++
			w.i = i
			b.WriteByte(' ')
			b.WriteString(w.writeval("cyan"))
		case constants.I_IncBy:
			b.WriteString(w.writeinst("INCB", "green"))
			i++
			w.i = i
			b.WriteByte(' ')
			b.WriteString(w.writeval("green"))
		case constants.I_DecBy:
			b.WriteString(w.writeinst("DECB", "green"))
			i++
			w.i = i
			b.WriteByte(' ')
			b.WriteString(w.writeval("green"))
		case constants.I_IncPBy:
			b.WriteString(w.writeinst("INCPB", "blue"))
			i++
			w.i = i
			b.WriteByte(' ')
			b.WriteString(w.writeval("blue"))
		case constants.I_DecPBy:
			b.WriteString(w.writeinst("DECPB", "blue"))
			i++
			w.i = i
			b.WriteByte(' ')
			b.WriteString(w.writeval("blue"))
		}
	}

	b.WriteByte(']')

	return b.String()
}
