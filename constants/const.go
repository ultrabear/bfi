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
	I_Zero uint = iota
	I_Inc
	I_Dec
	I_IncP
	I_DecP
	I_Read
	I_Write
	I_LStart
	I_LEnd
	I_IncBy
	I_DecBy
	I_IncPBy
	I_DecPBy
)
