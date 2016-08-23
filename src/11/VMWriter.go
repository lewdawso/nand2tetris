package vmwriter

import (
	"fmt"
	"os"
	"strconv"
)

type Segment int
type Command int

var f *os.File
var SegmentLookup map[Segment]string
var CommandLookup map[Command]string

const (
	CONST Segment = iota
	ARG
	LOCAL
	STATIC
	THIS
	THAT
	POINTER
	TEMP
)

const (
	ADD Command = iota
	SUB
	NEG
	EQ
	GT
	LT
	AND
	OR
	NOT
)

func init() {

	SegmentLookup = make(map[Segment]string)
	CommandLookup = make(map[Command]string)

	SegmentLookup[CONST] = "constant"
	SegmentLookup[ARG] = "arg"
	SegmentLookup[LOCAL] = "local"
	SegmentLookup[STATIC] = "static"
	SegmentLookup[THIS] = "this"
	SegmentLookup[THAT] = "that"
	SegmentLookup[POINTER] = "pointer"
	SegmentLookup[TEMP] = "temp"

	CommandLookup[ADD] = "add"
	CommandLookup[SUB] = "sub"
	CommandLookup[NEG] = "neg"
	CommandLookup[EQ] = "eq"
	CommandLookup[GT] = "gt"
	CommandLookup[LT] = "lt"
	CommandLookup[AND] = "and"
	CommandLookup[OR] = "or"
	CommandLookup[NOT] = "not"
}

//create output .vm file
func CreateFile(target string) {
	f, _ = os.Create(target)
}

func WritePush(s Segment, i int) {

	if str := SegmentLookup[s]; str == "" {
		fmt.Println("invalid segment")
		os.Exit(1)
	} else {
		write("push", str, strconv.Itoa(i))
	}
}

func WritePop(s Segment, i int) {

	if str := SegmentLookup[s]; str == "" {
		fmt.Println("invalid segment")
		os.Exit(1)
	} else {
		write("pop", str, strconv.Itoa(i))
	}
}

func WriteArithmetic(c Command) {

	if str := CommandLookup[c]; str == "" {
		fmt.Println("invalid command")
		os.Exit(1)
	} else {
		write(str, "", "")
	}
}

func WriteLabel(label string) {
	write("label", label, "")
}

func WriteGoto(label string) {
	write("goto", label, "")
}

func WriteIf(label string) {
	write("if-goto", label, "")
}

func WriteCall(name string, n int) {
	write("call", name, strconv.Itoa(n))
}

func WriteFunction(name string, n int) {
	write("function", name, strconv.Itoa(n))
}

func WriteReturn() {
	write("return", "", "")
}

func write(cmd string, arg1 string, arg2 string) {
	f.WriteString(cmd + " ")
	f.WriteString(arg1 + " ")
	f.WriteString(arg2 + " ")
	f.WriteString("\n")
}
