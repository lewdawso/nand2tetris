// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[3], respectively.)

// Put your code here.

// Calculate R0 * R1

@R0
D=M
@R3
M=D //Set R3 to R0; can then initiaite a decrementor without touching the actual data
@R2
M=0 //Reset R2 (product) to zero

(LOOP)
	@R3
	D=M
	@END
	D;JLE //if R3 has reached zero, branch to END
	@R1
	D=M
	@R2
	M=M+D //increase the product register by R1
	@R3
	M=M-1 //decrement "counter" by 1
	@LOOP
	0;JEQ
(END)
	@END
	0;JEQ
