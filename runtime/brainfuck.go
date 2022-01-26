package runtime

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/ultrabear/bfi/constants"
)

func indexByte(slice []byte, i int) *byte {
	return (*byte)(unsafe.Pointer(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&slice))) + uintptr(i)))
}

// Brainfuck represents a state of a brainfuck program and provides methods to run brainfuck code
type Brainfuck struct {
	pointer int
	buffer  []byte
	stdin   io.Reader
	stdout  io.Writer
}

// Inc increments the current cell
func (bfc *Brainfuck) Inc() {
	bfc.buffer[bfc.pointer]++
}

// IncUnsafe increments the current cell without bounds checking
func (bfc *Brainfuck) IncUnsafe() {
	(*indexByte(bfc.buffer, bfc.pointer))++
}

// Dec decrements the current cell
func (bfc *Brainfuck) Dec() {
	bfc.buffer[bfc.pointer]--
}

// DecUnsafe decrements the current cell without bounds checking
func (bfc *Brainfuck) DecUnsafe() {
	(*indexByte(bfc.buffer, bfc.pointer))--
}

// IncBy increments the current cell by a value
func (bfc *Brainfuck) IncBy(val uint) {
	bfc.buffer[bfc.pointer] += byte(val)
}

// IncByUnsafe increments the current cell by a value without bounds checking
func (bfc *Brainfuck) IncByUnsafe(val uint) {
	(*indexByte(bfc.buffer, bfc.pointer)) += byte(val)
}

// DecBy decrements the current cell by a value
func (bfc *Brainfuck) DecBy(val uint) {
	bfc.buffer[bfc.pointer] -= byte(val)
}

// DecByUnsafe decrements the current cell by a value without bounds checking
func (bfc *Brainfuck) DecByUnsafe(val uint) {
	(*indexByte(bfc.buffer, bfc.pointer)) -= byte(val)
}

// IncP increments the current pointer
func (bfc *Brainfuck) IncP() {
	bfc.pointer++
	if bfc.pointer >= len(bfc.buffer) {
		fmt.Fprintln(os.Stderr, constants.RuntimeOverflow)
		os.Exit(1)
	}
}

// DecP decrements the current pointer
func (bfc *Brainfuck) DecP() {
	bfc.pointer--
	if bfc.pointer < 0 {
		fmt.Fprintln(os.Stderr, constants.RuntimeUnderflow)
		os.Exit(1)
	}
}

// IncPBy increments the current pointer by a value
func (bfc *Brainfuck) IncPBy(amt uint) {
	bfc.pointer += int(amt)
	if bfc.pointer >= len(bfc.buffer) {
		fmt.Fprintln(os.Stderr, constants.RuntimeOverflow)
		os.Exit(1)
	}
}

// DecPBy decrements the current pointer by a value
func (bfc *Brainfuck) DecPBy(amt uint) {
	bfc.pointer -= int(amt)
	if bfc.pointer < 0 {
		fmt.Fprintln(os.Stderr, constants.RuntimeUnderflow)
		os.Exit(1)
	}
}

// Write writes the current cell to the stdout
func (bfc *Brainfuck) Write() {
	_, _ = bfc.stdout.Write(bfc.buffer[bfc.pointer : bfc.pointer+1])
}

// Read reads from stdin to the current cell
func (bfc *Brainfuck) Read() {

	n, err := bfc.stdin.Read(bfc.buffer[bfc.pointer : bfc.pointer+1])

	if err != nil || n == 0 {
		bfc.buffer[bfc.pointer] = 0
	}
}

// Zero will zero the current cell
func (bfc *Brainfuck) Zero() {
	bfc.buffer[bfc.pointer] = 0
}

// ZeroUnsafe will zero the current cell without bounds checking
func (bfc *Brainfuck) ZeroUnsafe() {
	(*indexByte(bfc.buffer, bfc.pointer)) = 0
}

// Cur will return the value of the current cell
func (bfc *Brainfuck) Cur() int {
	return int(bfc.buffer[bfc.pointer])
}

// CurUnsafe will return the value of the current cell without bounds checking
func (bfc *Brainfuck) CurUnsafe() int {
	return int(*indexByte(bfc.buffer, bfc.pointer))
}

// Initbfc will return a new Brainfuck instance with a passed buffer size
func Initbfc(size int) Brainfuck {
	bfc := Brainfuck{
		buffer:  make([]byte, size),
		pointer: 0,
		stdin:   os.Stdin,
		stdout:  os.Stdout,
	}

	return bfc

}
