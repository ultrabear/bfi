# bfi
A commandline brainfuck interpreter written to balance execution speed and convenience
# Basically
I had wanted a commandline tool to execute brainfuck inline with scripts, why did I want that? dunno but it stuck  
The initial version of bfi was written in python, and then ported to golang where it got actual care to running fast  
While the goal of bfi is to be scriptable and convenient to just boot up on the commandline, it also hopes to be fast enough to handle most brainfuck in a timely matter, with a static compiling stage that doesnt run any brainfuck code while trying to optimize it  
# Installation/Usage
Install
```
sudo make install
```
Usage (argv)
```
bfi "++++[>++++[>++++<-]>>++<<<-]>>+.>++." # Prints A and a newline
```
Usage (files)
```
bfi f <filename>
```
# This implementation
This implementation has 2 notable differences from most brainfuck compilers/interpreters
- It will compress opposing instructions out of the loop entirely
  - What this means is that `+-` will be compiled to `NIL` and `<>` will be compiled to `NIL`
  - A special behaviour of this is that `<` will raise an error for underflowing the pointer location but `<>` will raise no error as before it gets to runtime that is compiled out
- The buffer size is calculated at compile time based off the size of the input stream + 1
  - This will in theory never error on code that is stable as you can only move forwards once per instruction, but code that "flies" on its own such as `+[>+]` will almost instantly hit the end of the buffer and error  
Other notes on implementation  
- The execution environment is 8 bit brainfuck, there are no settings to use other sizes  
 
