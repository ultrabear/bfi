package runtime

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/ultrabear/bfi/constants"
)

func IndexByte(slice []byte, i int) *byte {
	return (*byte)(unsafe.Pointer(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&slice))) + uintptr(i)))
}

type Brainfuck struct {
	pointer int
	buffer  []byte
	stdin   io.Reader
	stdout  io.Writer
}

func (bfc *Brainfuck) Inc() {
	bfc.buffer[bfc.pointer]++
}

func (bfc *Brainfuck) IncUnsafe() {
	(*IndexByte(bfc.buffer, bfc.pointer))++
}

func (bfc *Brainfuck) Dec() {
	bfc.buffer[bfc.pointer]--
}

func (bfc *Brainfuck) DecUnsafe() {
	(*IndexByte(bfc.buffer, bfc.pointer))--
}

func (bfc *Brainfuck) IncBy(val uint) {
	bfc.buffer[bfc.pointer] += byte(val)
}

func (bfc *Brainfuck) IncByUnsafe(val uint) {
	(*IndexByte(bfc.buffer, bfc.pointer)) += byte(val)
}

func (bfc *Brainfuck) DecBy(val uint) {
	bfc.buffer[bfc.pointer] -= byte(val)
}

func (bfc *Brainfuck) DecByUnsafe(val uint) {
	(*IndexByte(bfc.buffer, bfc.pointer)) -= byte(val)
}

func (bfc *Brainfuck) IncP() {
	bfc.pointer++
	if bfc.pointer >= len(bfc.buffer) {
		fmt.Fprintln(os.Stderr, constants.RuntimeOverflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) DecP() {
	bfc.pointer--
	if bfc.pointer < 0 {
		fmt.Fprintln(os.Stderr, constants.RuntimeUnderflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) IncPBy(amt uint) {
	bfc.pointer += int(amt)
	if bfc.pointer >= len(bfc.buffer) {
		fmt.Fprintln(os.Stderr, constants.RuntimeOverflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) DecPBy(amt uint) {
	bfc.pointer -= int(amt)
	if bfc.pointer < 0 {
		fmt.Fprintln(os.Stderr, constants.RuntimeUnderflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) Write() {
	_, _ = bfc.stdout.Write(bfc.buffer[bfc.pointer : bfc.pointer+1])
}

func (bfc *Brainfuck) Read() {

	n, err := bfc.stdin.Read(bfc.buffer[bfc.pointer : bfc.pointer+1])

	if err != nil || n == 0 {
		bfc.buffer[bfc.pointer] = 0
	}
}

func (bfc *Brainfuck) Zero() {
	bfc.buffer[bfc.pointer] = 0
}

func (bfc *Brainfuck) ZeroUnsafe() {
	(*IndexByte(bfc.buffer, bfc.pointer)) = 0
}

func (bfc *Brainfuck) Cur() int {
	return int(bfc.buffer[bfc.pointer])
}

func (bfc *Brainfuck) CurUnsafe() int {
	return int(*IndexByte(bfc.buffer, bfc.pointer))
}

func Initbfc(size int) Brainfuck {
	bfc := Brainfuck{
		buffer:  make([]byte, size),
		pointer: 0,
		stdin:   os.Stdin,
		stdout:  os.Stdout,
	}

	return bfc

}
