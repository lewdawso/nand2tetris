/*
Symbol table abstraction. Associates identifier names found in the program
with identifier properties needed for compilation.
*/

package symtable

import (
	"fmt"
)

var classSymbolTable map[string]Symbol
var subroutineSymbolTable map[string]Symbol
var indices map[Kind]int

type Symbol struct {
	_type string
	kind  Kind
	index int
}

type Kind int

const (
	NONE Kind = iota
	STATIC
	FIELD
	ARG
	VAR
)

func init() {
	classSymbolTable = make(map[string]Symbol)
	subroutineSymbolTable = make(map[string]Symbol)
	indices := make(map[Kind]int)

	indices[STATIC] = 0
	indices[FIELD] = 0
	indices[ARG] = 0
	indices[VAR] = 0
}

func startSubroutine() {
	subroutineSymbolTable = make(map[string]Symbol)
}

func Define(name string, _type string, kind Kind) {
	//find the count of this kind of variable
	index := VarCount(kind) + 1
	indices[kind] = index

	//define Symbol variable
	symbol := Symbol{_type, kind, index}

	if kind == STATIC || kind == FIELD {
		classSymbolTable[name] = symbol
	} else if kind == ARG || kind == VAR {
		subroutineSymbolTable[name] = symbol
	} else {
		fmt.Println("unrecognised kind")
	}
}

func VarCount(kind Kind) int {
	return indices[kind]
}

func KindOf(name string) Kind {
	s := classSymbolTable[name]
	if !EmptySymbol(s) {
		return s.kind
	} else {
		return subroutineSymbolTable[name].kind
	}
}

func TypeOf(name string) string {
	if s := classSymbolTable[name]; s._type != "" {
		return s._type
	} else {
		return subroutineSymbolTable[name]._type
	}
}

func IndexOf(name string) Symbol {
	if s := classSymbolTable[name]; s.index == 0 {
		return s
	} else {
		return subroutineSymbolTable[name]
	}
}

func EmptySymbol(s Symbol) bool {
	if s.kind == 0 && s._type == "" && s.index == 0 {
		return true
	}
	return false
}
