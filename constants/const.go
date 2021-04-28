package constants

// Define errors here to avoid looking for them in the code
const (
	Error                string = "\033[91mERROR:\033[0m "
	SyntaxEndBeforeStart string = Error + "BF Syntax: Loop end defined before loop start"
	SyntaxUnbalanced     string = Error + "BF Syntax: Unbalanced loop statements"
	RuntimeUnderflow     string = Error + "BF Runtime: Underflowed pointer location"
	RuntimeOverflow      string = Error + "BF Runtime: Overflowed pointer location"
)
