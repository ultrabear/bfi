# ultrabear/bfi
A commandline brainfuck interpreter written to balance execution speed and convenience
# Basically
I had wanted a commandline tool to execute brainfuck inline with scripts, why did I want that? dunno but it stuck  
The initial version of bfi was written in python, and then ported to golang where it got actual care to running fast  
While the goal of bfi is to be scriptable and convenient to just boot up on the commandline, it also hopes to be fast enough to handle most brainfuck in a timely matter, with a static compiling stage that doesnt run any brainfuck code while trying to optimize it  
# Installation/Usage
Install
```bash
sudo make install
```
Usage (argv)
```bash
bfi "++++[>++++[>++++<-]>>++<<<-]>>+.>++." # Prints A and a newline
```
Usage (files)
```bash
bfi f <filename>
```
# This implementation
## Notable Differences
- It will compress opposing instructions out of the loop entirely
  - What this means is that `+-` will be compiled to `NIL` and `<>` will be compiled to `NIL`
  - A special behaviour of this is that `<` will raise an error for underflowing the pointer location but `<>` will raise no error as before it gets to runtime that is compiled out
- The buffer size is calculated at compile time based off the count of `>` instructions or 30k as a minimum, for the brainfuck standard
  - In most cases this will not error on brainfuck, but some programs rely on buffer flying and may not function correctly
- `[]` is compiled out entirely
  - This was in a weak attempt to strip some comments from the program so that it does not affect runtime
	- Notably programs such as `+[]` will not loop forever and will instead just run as `+` on another interpreter
	- As `[]` is essentially an infinite loop this will not affect most programs
## Other Notes  
- The execution environment is 8 bit brainfuck, there are no settings to use other sizes  
# Credits
All example programs (except beemovie.bf) came from the example sets from [El Brainfuck](https://copy.sh/brainfuck), a bf-to-JS compilier implementation with debugging tools and variable buffer widths
