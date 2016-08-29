package main

import (
	"fmt"
	"io/ioutil"
	"nand2tetris/SymTable"
	"nand2tetris/VMWriter"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var tokens [][]string
var current []string
var output []string

var operators = []string{"+", "-", "*", "/", "&amp;", "|", "&lt;", "&gt;", "="}
var unaryOp = []string{"-", "~"}

var currentClass string
var currentSubroutine string

var kind symtable.Kind
var _type string

var labelCount int

func generateTokenArray(slice []string) {
	re, _ := regexp.Compile(" ")
	for i := range slice {
		pair := re.Split(slice[i], 2)
		tokens = append(tokens, pair)
	}
	advance()
}

func advance() {
	current = tokens[0]
	tokens = tokens[1:]
}

func raiseError(message string) {
	fmt.Println("error: ", message)
}

func checkToken(token string) bool {
	if strings.Compare(token, current[1]) == 0 {
		return true
	}
	return false
}

func checkTokenSlice(slice []string) bool {
	for i := range slice {
		if strings.Compare(slice[i], current[1]) == 0 {
			return true
		}
	}
	return false
}

func checkTokenTypeSlice(slice []string) bool {
	for i := range slice {
		if strings.Compare(slice[i], current[0]) == 0 {
			return true
		}
	}
	return false
}

func checkTokenTypeSliceCustom(_type string, slice []string) bool {
	for i := range slice {
		if strings.Compare(slice[i], _type) == 0 {
			return true
		}
	}
	return false
}

func checkIdentifier() bool {
	re, err := regexp.Compile("^[a-zA-Z_][\\w]*$")
	if err != nil {
		os.Exit(1)
	}
	if !re.MatchString(current[1]) {
		raiseError("invalid identifier")
		return false
	}
	return true
}

func checkIdentifierPassive() bool {
	re, err := regexp.Compile("^[a-zA-Z_][\\w]*$")
	if err != nil {
		os.Exit(1)
	}
	if !re.MatchString(current[1]) {
		return false
	}
	return true
}

func checkOpeningBrace() bool {
	if !checkToken("{") {
		raiseError("missing opening brace")
		return false
	}
	advance()
	return true
}

func checkClosingBrace() bool {
	if !checkToken("}") {
		raiseError("missing closing brace")
		return false
	}
	advance()
	return true
}

func checkOpeningBracket() bool {
	if !checkToken("(") {
		raiseError("missing opening bracket")
		return false
	}
	advance()
	return true
}

func checkClosingBracket() bool {
	if !checkToken(")") {
		raiseError("missing closing bracket")
		return false
	}
	advance()
	return true
}

func checkSemicolon() bool {
	if !checkToken(";") {
		raiseError("missing semicolon")
		return false
	}
	advance()
	return true
}

func checkTypeAndIdentifier() bool {
	//check type
	if !checkTokenSlice([]string{"int", "char", "bool"}) && !checkIdentifierPassive() {
		raiseError("type specifier")
		return false
	}

	_type = getCurrent()
	advance()

	if !checkIdentifier() {
		return false
	}

	name := getCurrent()
	symtable.Define(name, _type, kind)
	advance()

	return true
}

func checkNextToken(token string) bool {
	if strings.Compare(token, tokens[0][1]) == 0 {
		return true
	}
	return false
}

func debug() {
	fmt.Println("current: ", current[1])
	fmt.Println("next: ", tokens[0][1])

	for i := range output {
		fmt.Println(output[i])
	}
}

func getCurrent() string {
	return current[1]
}

func getCurrentInt() int {
	n, err := strconv.Atoi(current[1])
	if err != nil {
		os.Exit(1)
	}
	return n
}

func getSegment(kind symtable.Kind) vmwriter.Segment {
	switch kind {
	case symtable.STATIC:
		return vmwriter.STATIC
	case symtable.FIELD:
		return vmwriter.THIS
	case symtable.ARG:
		return vmwriter.ARG
	case symtable.VAR:
		return vmwriter.LOCAL
	}
	//FIX
	return vmwriter.LOCAL
}

func genLabel() string {
	count := labelCount
	labelCount++
	return "LABEL" + strconv.Itoa(count)
}

func compileClass() bool {

	if !checkToken("class") {
		raiseError("missing class keyword to open")
		return false
	}

	advance()

	if !checkIdentifier() {
		return false
	}

	currentClass = getCurrent()
	advance()

	if !checkOpeningBrace() {
		return false
	}

	for checkTokenSlice([]string{"static", "field"}) {
		if !compileClassVarDec() {
			raiseError("compileClassVarDec")
			return false
		}
	}

	for checkTokenSlice([]string{"constructor", "function", "method"}) {
		if !compileSubroutine() {
			raiseError("compileSubroutine")
			return false
		}
	}

	if !checkClosingBrace() {
		return false
	}

	if len(tokens) != 0 {
		raiseError("tokens remaining after class closed")
		return false
	}

	return true
}

func compileClassVarDec() bool {

	kind = symtable.KindLookup[getCurrent()]

	advance()

	//check type
	if !checkTokenSlice([]string{"int, char, bool"}) && !checkIdentifierPassive() {
		raiseError("missing type specifier")
		return false
	}

	_type := getCurrent()
	advance()

	if !checkIdentifier() {
		return false
	}

	name := getCurrent()
	advance()

	symtable.Define(name, _type, kind)

	//deal with (',', varName)*
	for checkToken(",") {
		advance()
		if !checkIdentifier() {
			return false
		}
		name = getCurrent()
		symtable.Define(name, _type, kind)
		advance()
	}

	if !checkSemicolon() {
		return false
	}

	return true
}

func compileSubroutine() bool {

	//need to log whether we're dealing with a method or a function
	method := false

	symtable.StartSubroutine()

	//first argument of a method is the object itself
	if getCurrent() == "method" {
		symtable.Define("this", currentClass, symtable.ARG)
		method = true
	}

	advance()

	//check return type
	if !checkTokenSlice([]string{"void", "int", "char", "bool"}) && !checkIdentifierPassive() {
		raiseError("missing return type")
		return false
	}

	advance()

	//check subroutine name
	if !checkIdentifier() {
		raiseError("invalid identifier")
		return false
	}

	currentSubroutine = getCurrent()
	advance()

	if !checkOpeningBracket() {
		return false
	}

	kind = symtable.ARG
	if !compileParameterList() {
		raiseError("compileParameterList")
		return false
	}

	if !checkClosingBracket() {
		return false
	}

	if !checkOpeningBrace() {
		return false
	}

	//possible variable declaration
	kind = symtable.VAR
	for checkToken("var") {
		if !compileVarDec() {
			raiseError("compileVarDec")
			return false
		}
	}

	vmwriter.WriteFunction(currentClass+"."+currentSubroutine, symtable.VarCount(symtable.VAR))
	//if this is a method, need to set the "this" pointer
	if method {
		vmwriter.WritePush(vmwriter.ARG, 0)
		vmwriter.WritePop(vmwriter.POINTER, 0)
	}

	//compile statements
	if !compileStatements() {
		raiseError("compileStatements")
		return false
	}

	if !checkClosingBrace() {
		return false
	}

	return true
}

func compileParameterList() bool {

	//no parameters
	if checkToken(")") {
		return true
	}

	if !checkTypeAndIdentifier() {
		raiseError("checkTypeAndIdentifier")
		return false
	}

	for checkToken(",") {
		advance()
		if !checkTypeAndIdentifier() {
			raiseError("checkTypeAndIdentifier")
			return false
		}
	}

	return true
}

func compileVarDec() bool {
	advance()

	if !checkTypeAndIdentifier() {
		raiseError("checkTypeAndIdentifier")
		return false
	}

	for checkToken(",") {
		advance()
		if !checkIdentifier() {
			return false
		}
		name := getCurrent()
		symtable.Define(name, _type, kind)
		advance()
	}

	if !checkSemicolon() {
		return false
	}

	return true
}

func compileStatements() bool {
	for checkTokenSlice([]string{"let", "if", "while", "do", "return"}) {

		switch current[1] {
		case "let":
			if !compileLet() {
				raiseError("compileLet")
				return false
			}
			break
		case "if":
			if !compileIf() {
				raiseError("compileIf")
				return false
			}
			break
		case "while":
			if !compileWhile() {
				raiseError("compileWhile")
				return false
			}
			break
		case "do":
			if !compileDo() {
				raiseError("compileDo")
				return false
			}
			break
		case "return":
			if !compileReturn() {
				raiseError("compileReturn")
				return false
			}
			break
		default:
			raiseError("invalid statement keyword")
			return false
		}
	}
	return true
}

func compileLet() bool {
	offset := false
	advance()

	if !checkIdentifier() {
		return false
	}

	varName := getCurrent()
	advance()

	//potential expression
	if checkToken("[") {
		offset = true

		//push array variable onto the stack (its base address)
		vmwriter.WritePush(getSegment(symtable.KindOf(varName)), symtable.IndexOf(varName))
		advance()

		if !compileExpression() {
			raiseError("compileExpression")
			return false
		}

		if !checkToken("]") {
			raiseError("missing closing ]")
			return false
		}
		advance()

		//add output of expression to base address which was pushed onto the
		//stack

		vmwriter.WriteArithmetic(vmwriter.ADD)

	}

	//equals expression
	if !checkToken("=") {
		raiseError("missing =")
		return false
	}
	advance()

	if !compileExpression() {
		raiseError("compileExpression")
		return false
	}

	if !checkSemicolon() {
		return false
	}

	if offset {
		//stack looks like this:
		// base + offset
		// expression result

		//store expression result in TEMP
		vmwriter.WritePop(vmwriter.TEMP, 0)
		//align THAT with (base + offset)
		vmwriter.WritePop(vmwriter.POINTER, 1)
		//push result to stack
		vmwriter.WritePush(vmwriter.TEMP, 0)
		//pop result to THAT
		vmwriter.WritePop(vmwriter.THAT, 0)
	} else {
		//stack looks like this:
		//expression result
		vmwriter.WritePop(getSegment(symtable.KindOf(varName)), symtable.IndexOf(varName))
	}
	return true
}

func compileIf() bool {
	advance()
	start := genLabel()
	end := genLabel()

	if !checkOpeningBracket() {
		return false
	}
	if !compileExpression() {
		raiseError("compileExpression")
		return false
	}
	if !checkClosingBracket() {
		return false
	}

	vmwriter.WriteArithmetic(vmwriter.NOT)
	vmwriter.WriteIf(start)

	if !checkOpeningBrace() {
		return false
	}
	if !compileStatements() {
		raiseError("compileStatements")
		return false
	}
	if !checkClosingBrace() {
		return false
	}

	vmwriter.WriteGoto(end)
	vmwriter.WriteLabel(start)

	if checkToken("else") {
		advance()
		if !checkOpeningBrace() {
			return false
		}
		if !compileStatements() {
			return false
		}
		if !checkClosingBrace() {
			return false
		}
	}
	vmwriter.WriteLabel(end)
	return true
}

func compileWhile() bool {
	advance()
	loop := genLabel()
	exit := genLabel()

	vmwriter.WriteLabel(loop)

	if !checkOpeningBracket() {
		return false
	}
	if !compileExpression() {
		raiseError("compileExpression")
		return false
	}
	if !checkClosingBracket() {
		return false
	}

	vmwriter.WriteArithmetic(vmwriter.NOT)
	vmwriter.WriteIf(exit)

	if !checkOpeningBrace() {
		return false
	}
	if !compileStatements() {
		raiseError("compileStatements")
		return false
	}
	if !checkClosingBrace() {
		return false
	}

	vmwriter.WriteGoto(loop)
	vmwriter.WriteLabel(exit)

	return true
}

func compileDo() bool {
	advance()

	if !checkSubroutineCall() {
		raiseError("subroutineCall")
		return false
	}
	if !checkSemicolon() {
		return false
	}

	return true
}

func compileReturn() bool {
	advance()

	//expression is optional
	if checkToken(";") {
		vmwriter.WritePush(vmwriter.CONST, 0)
		vmwriter.WriteReturn()
		advance()
		return true
	}

	if !compileExpression() {
		raiseError("compileExpression")
		return true
	}
	if !checkSemicolon() {
		return false
	}

	vmwriter.WriteReturn()
	return true
}

func compileExpression() bool {
	//an expression must contain at least one term
	if !compileTerm() {
		raiseError("compileTerm")
		return false
	}
	//(op term)*
	op := true
	for op {
		for i := range operators {
			if checkToken(operators[i]) {
				op = true
				//store the operator until the next term is pushed to the stack
				operator := getCurrent()
				advance()
				if !compileTerm() {
					raiseError("compileTerm")
					return false
				}
				//now write the operator code
				switch operator {
				case "+":
					vmwriter.WriteArithmetic(vmwriter.ADD)
				case "-":
					vmwriter.WriteArithmetic(vmwriter.SUB)
				case "*":
					vmwriter.WriteCall("Math.multiply", 2)
				case "/":
					vmwriter.WriteCall("Math.divide", 2)
				case "&amp;":
					vmwriter.WriteArithmetic(vmwriter.AND)
				case "|":
					vmwriter.WriteArithmetic(vmwriter.OR)
				case "&lt;":
					vmwriter.WriteArithmetic(vmwriter.LT)
				case "&gt;":
					vmwriter.WriteArithmetic(vmwriter.GT)
				case "=":
					vmwriter.WriteArithmetic(vmwriter.EQ)
				default:
					raiseError("unknown operator")
				}
				//break out of for loop
				break
			}
			op = false
		}
	}

	return true
}

func compileTerm() bool {
	//int, string or keyword
	if checkTokenTypeSlice([]string{"integerConstant", "stringConstant", "keywordConstant"}) {
		switch current[0] {
		case "integerConstant":
			vmwriter.WritePush(vmwriter.CONST, getCurrentInt())
		case "stringConstant":
			str := getCurrent()
			//push length of string onto stack
			vmwriter.WritePush(vmwriter.CONST, len(str))
			//call String.new() which has 1 argument
			vmwriter.WriteCall("String.new", 1)
			for _, char := range str {
				vmwriter.WritePush(vmwriter.CONST, int(char))
				//2 arguments because first argument is the object itself (this)
				//which is left on the stack after the call to String.new
				vmwriter.WriteCall("String.appendChar", 2)
			}
		case "keywordConstant":
			if current[1] == "true" {
				vmwriter.WritePush(vmwriter.CONST, 1)
				vmwriter.WriteArithmetic(vmwriter.NOT)
			} else if current[1] == "false" {
				vmwriter.WritePush(vmwriter.CONST, 0)
			} else if current[1] == "null" {
				vmwriter.WritePush(vmwriter.CONST, 0)
			} else {
				raiseError("invalid keywordConstant value")
				return false
			}
		default:
			raiseError("unknown type qualifier")
			return false
		}
		advance()
		return true
	}

	//(expression)
	if checkToken("(") {
		advance()
		if !compileExpression() {
			raiseError("compileExpression")
			return false
		}
		if !checkClosingBracket() {
			return false
		}
		return true
	}

	//unaryOp term
	for i := range unaryOp {
		if strings.Compare(unaryOp[i], current[1]) == 0 {
			var op string = getCurrent()

			advance()
			if !compileTerm() {
				raiseError("compileTerm")
				return false
			}
			switch op {
			case "-":
				vmwriter.WriteArithmetic(vmwriter.NEG)
			case "~":
				vmwriter.WriteArithmetic(vmwriter.NOT)
			default:
				raiseError("unknown unary operator")
				return false
			}
			return true
		}
	}

	//can now only have varName | varName [expression] | subroutineCall
	//All of these terms begin with varName => look ahead one token to
	//differentiate

	if checkNextToken("[") {
		if !checkIdentifier() {
			return false
		}

		varName := getCurrent()
		vmwriter.WritePush(getSegment(symtable.KindOf(varName)), symtable.IndexOf(varName))
		advance()
		advance()

		if !compileExpression() {
			raiseError("compileExpression")
			return false
		}

		if !checkToken("]") {
			raiseError("missing ]")
			return false
		}
		advance()

		//stack now looks like this
		//varName
		//expression result

		//ADD
		vmwriter.WriteArithmetic(vmwriter.ADD)
		//set THAT to varName + offset
		vmwriter.WritePop(vmwriter.POINTER, 1)
		//push value to the stack
		vmwriter.WritePush(vmwriter.THAT, 0)
		return true
	}

	if checkNextToken("(") || checkNextToken(".") {
		if !checkSubroutineCall() {
			raiseError("checkSubroutineCall")
			return false
		}
		return true
	}

	if !checkIdentifier() {
		return false
	}

	vmwriter.WritePush(getSegment(symtable.KindOf(getCurrent())), symtable.IndexOf(getCurrent()))
	advance()

	return true
}

func checkSubroutineCall() bool {
	if !checkIdentifier() {
		return false
	}

	var args, additional int = 0, 0
	var name string
	subroutine := getCurrent()
	advance()

	//could be subroutineName(blargh) | class/var.subroutine(blargh)

	if checkToken(".") {
		advance()
		if !checkIdentifier() {
			return false
		}
		object := subroutine
		subroutine := getCurrent()

		//"object" could be a class or a variable (instance of a class)
		_type := symtable.TypeOf(object)
		if checkTokenTypeSliceCustom(_type, []string{"int", "boolean", "char", "void"}) {
			raiseError("not a valid object type")
		} else if _type == "" {
			//it's a class
			name = object + "." + subroutine
		} else {
			//it's an instance of a class
			args++
			vmwriter.WritePush(getSegment(symtable.KindOf(object)), symtable.IndexOf(object))
			name = _type + "." + subroutine
		}
	} else {
		name = currentClass + "." + subroutine
		args++
		vmwriter.WritePush(vmwriter.POINTER, 0)
	}
	advance()

	if !checkOpeningBracket() {
		return false
	}

	if additional = compileExpressionList(); additional == -1 {
		raiseError("compileExpressionList")
		return false
	}
	if !checkClosingBracket() {
		return false
	}

	args += additional
	vmwriter.WriteCall(name, args)
	return true
}

func compileExpressionList() int {
	var args int

	//first off, check if we have an empty list
	if checkToken(")") {
		return 0
	}

	//now we know we have to compile at least one expression
	if !compileExpression() {
		raiseError("compileExpression")
		return -1
	}

	args++

	for !checkToken(")") {
		if !checkToken(",") {
			raiseError("missing comma")
			return -1
		}
		advance()
		if !compileExpression() {
			raiseError("compileExpression")
			return -1
		}
		args++
	}
	return args
}

func main() {
	//arguments
	args := os.Args
	//tokens
	filepath := args[1]
	//output file
	target := args[2]

	//open tokens file
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		os.Exit(1)
	}

	stringified := string(data)
	re, _ := regexp.Compile("\"")
	stringified = re.ReplaceAllString(stringified, "")
	slice := strings.Split(stringified, "\n")
	generateTokenArray(slice)

	//sort op slice
	sort.Strings(operators)

	//create output file
	vmwriter.CreateFile(target)

	//first routine to be called must be compileClass
	if !compileClass() {
		raiseError("unable to compile class")
		debug()
	}
}
