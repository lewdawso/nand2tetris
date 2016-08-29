/*
Symbol table abstraction. Associates identifier names found in the program
with identifier properties needed for compilation.
*/

package symtable

import (
	"fmt"
	"os"
)

var classSymbolTable map[string]Symbol
var subroutineSymbolTable map[string]Symbol
var indices map[Kind]int

var KindLookup map[string]Kind

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
	indices = make(map[Kind]int)
	KindLookup = make(map[string]Kind)

	indices[STATIC] = 0
	indices[FIELD] = 0
	indices[ARG] = 0
	indices[VAR] = 0

	KindLookup["static"] = STATIC
	KindLookup["field"] = FIELD
	KindLookup["arg"] = ARG
	KindLookup["var"] = VAR
}

func StartSubroutine() {
	subroutineSymbolTable = make(map[string]Symbol)
	indices[ARG] = 0
	indices[VAR] = 0
}

func Define(name string, _type string, kind Kind) {
	//find the count of this kind of variable
	index := VarCount(kind)
	indices[kind] = index + 1

	//define Symbol variable
	symbol := Symbol{_type, kind, index}

	if kind == STATIC || kind == FIELD {
		classSymbolTable[name] = symbol
	} else if kind == ARG || kind == VAR {
		subroutineSymbolTable[name] = symbol
	} else {
		fmt.Println("unrecognised kind")
		os.Exit(1)
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

func IndexOf(name string) int {
	if s := classSymbolTable[name]; !EmptySymbol(s) {
		return s.index
	} else if s := subroutineSymbolTable[name]; !EmptySymbol(s) {
		return s.index
	} else {
		return -1
	}
}

func EmptySymbol(s Symbol) bool {
	if s.kind == 0 && s._type == "" && s.index == 0 {
		return true
	}
	return false
}
