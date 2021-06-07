package runtime

import (
	"fmt"
	"github.com/ultrabear/bfi/constants"
	"unsafe"
	"io"
	"os"
)

type Brainfuck struct {
	unbuffer uintptr
	buffer  []byte
	pointer int
	stdin   io.Reader
	stdout  io.Writer
}

func (bfc *Brainfuck) Inc() {
	bfc.buffer[bfc.pointer]++
}

func (bfc *Brainfuck) IncUnsafe() {
	*(*byte)(unsafe.Pointer(bfc.unbuffer+uintptr(bfc.pointer))) += 1
}

func (bfc *Brainfuck) Dec() {
	bfc.buffer[bfc.pointer]--
}

func (bfc *Brainfuck) DecUnsafe() {
	*(*byte)(unsafe.Pointer(bfc.unbuffer+uintptr(bfc.pointer))) -= 1
}

func (bfc *Brainfuck) IncBy(val uint) {
	bfc.buffer[bfc.pointer] += byte(val)
}

func (bfc *Brainfuck) IncByUnsafe(val uint) {
	*(*byte)(unsafe.Pointer(bfc.unbuffer+uintptr(bfc.pointer))) += byte(val)
}

func (bfc *Brainfuck) DecBy(val uint) {
	bfc.buffer[bfc.pointer] -= byte(val)
}

func (bfc *Brainfuck) DecByUnsafe(val uint) {
	*(*byte)(unsafe.Pointer(bfc.unbuffer+uintptr(bfc.pointer))) -= byte(val)
}

func (bfc *Brainfuck) IncP() {
	bfc.pointer += 1
	if bfc.pointer >= len(bfc.buffer) {
		fmt.Println(constants.RuntimeOverflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) DecP() {
	bfc.pointer -= 1
	if bfc.pointer < 0 {
		fmt.Println(constants.RuntimeUnderflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) IncPBy(amt uint) {
	bfc.pointer += int(amt)
	if bfc.pointer >= len(bfc.buffer) {
		fmt.Println(constants.RuntimeOverflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) DecPBy(amt uint) {
	bfc.pointer -= int(amt)
	if bfc.pointer < 0 {
		fmt.Println(constants.RuntimeUnderflow)
		os.Exit(1)
	}
}

func (bfc *Brainfuck) Write() {
	bfc.stdout.Write([]byte{bfc.buffer[bfc.pointer]})
}

func (bfc *Brainfuck) Read() {
	var indata = make([]byte, 1)
	amt, err := bfc.stdin.Read(indata)
	if amt != 1 {
		indata[0] = 0
	}
	if err != nil {
		indata[0] = 0
	}
	bfc.buffer[bfc.pointer] = indata[0]
}

func (bfc *Brainfuck) Zero() {
	bfc.buffer[bfc.pointer] = 0
}

func (bfc *Brainfuck) ZeroUnsafe() {
	*(*byte)(unsafe.Pointer(bfc.unbuffer+uintptr(bfc.pointer))) = 0
}

func (bfc *Brainfuck) Cur() int {
	return int(bfc.buffer[bfc.pointer])
}

func (bfc *Brainfuck) CurUnsafe() int {
	return int(*(*byte)(unsafe.Pointer(bfc.unbuffer+uintptr(bfc.pointer))))
}

type FakeSlice struct {
	Pointer uintptr
	Len uintptr
	Cap uintptr
}

func Initbfc(size int) Brainfuck {
	bfc := Brainfuck{
		buffer:  make([]byte, size),
		pointer: 0,
		stdin:   os.Stdin,
		stdout:  os.Stdout,
	}

	bfc.unbuffer = ((*FakeSlice)(unsafe.Pointer(&bfc.buffer))).Pointer

	return bfc

}
