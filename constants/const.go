// Package constants contains predefined constants for bfc format
package constants

// Define errors here to avoid looking for them in the code
const (
	Error                = "\033[91mERROR:\033[0m "
	SyntaxEndBeforeStart = Error + "BF Syntax: Loop end defined before loop start"
	SyntaxUnbalanced     = Error + "BF Syntax: Unbalanced loop statements"
	RuntimeUnderflow     = Error + "BF Runtime: Underflowed pointer location"
	RuntimeOverflow      = Error + "BF Runtime: Overflowed pointer location"
)

// Define instructions
const (
	InstrucZero uint = iota
	InstrucInc
	InstrucDec
	InstrucIncP
	InstrucDecP
	InstrucRead
	InstrucWrite
	InstrucLStart
	InstrucLEnd
	InstrucIncBy
	InstrucDecBy
	InstrucIncPBy
	InstrucDecPBy
)
