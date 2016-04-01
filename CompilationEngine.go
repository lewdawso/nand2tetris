package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "strings"
    "regexp"
    "sort"
)

var tokens [][]string
var current []string
var output []string

var operators = []string{"+", "-", "*", "/", "&amp;", "|", "&lt;", "&gt;", "="}
var unaryOp = []string{"-", "~"}

func generateTokenArray(slice []string) {

    for i := range slice {
        pair := strings.Split(slice[i], " ")
        tokens = append(tokens, pair)
    }
    advance()
}

func write(cmd string) {
    output = append(output, cmd)
}

func writeOpen(cmd string) {
    output = append(output, "<" + cmd + ">")
}

func writeClose(cmd string) {
    output = append(output, "</" + cmd + ">")
}

func writeToken() {
    //do this better
    formatted := "<" + current[0] + ">" + " " + current[1] + " " + "</" + current[0] + ">"
    output = append(output, formatted)
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

func checkIdentifier() bool {
    re, err := regexp.Compile("^[a-zA-Z_][\\w]*$")
    if err != nil {
        os.Exit(1)
    }
    if !re.MatchString(current[1]) {
        raiseError("invalid identifier")
        return false
    }
    writeToken()
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
    writeToken()
    return true
}

func checkClosingBrace() bool {
    if !checkToken("}") {
        raiseError("missing closing brace")
        return false
    }
    writeToken()
    return true
}

func checkOpeningBracket() bool {
    if !checkToken("(") {
        raiseError("missing opening bracket")
        return false
    }
    writeToken()
    return true
}

func checkClosingBracket() bool {
    if !checkToken(")") {
        raiseError("missing closing bracket")
        return false
    }
    writeToken()
    return true
}

func checkSemicolon() bool {
    if !checkToken(";") {
        raiseError("missing semicolon")
        return false
    }
    writeToken()
    return true
}

func checkTypeAndIdentifier() bool {
    //check type
    if (!checkTokenSlice([]string{"int", "char", "bool"}) && !checkIdentifierPassive()) {
        raiseError("type specifier")
        return false
    }

    writeToken()

    if !checkIdentifier() {
        return false
    }
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

func compileClass() bool {

    if !checkToken("class") {
        raiseError("missing class keyword to open")
        return false
    }

    writeOpen("class")
    writeToken()

    if !checkIdentifier() { return false }

    if !checkOpeningBrace() { return false }

    for (checkTokenSlice([]string{"static", "field"})) {
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

    if !checkClosingBrace() { return false }

    if len(tokens) != 0 {
        raiseError("tokens remaining after class closed")
        return false
    }

    writeClose("class")
    return true
}

func compileClassVarDec() bool {
    writeOpen("classVarDec")
    writeToken()

    //check type
    if (!checkTokenSlice([]string{"int, char, bool"}) && !checkIdentifierPassive()) {
        raiseError("missing type specifier")
        return false
    }

    writeToken()

    if !checkIdentifier() { return false }

    //deal with (',', varName)*
    for checkToken(",") {
        writeToken()
        if !checkIdentifier() {
            return false
        }
    }

    if !checkSemicolon() { return false }

    writeClose("classVarDec")
    return true
}

func compileSubroutine() bool {
    writeOpen("subroutineDec")
    writeToken()

    //check return type
    if !checkTokenSlice([]string{"void", "int", "char", "bool"}) && !checkIdentifierPassive() {
        raiseError("missing return type")
        return false
    }

    writeToken()

    //check subroutine name
    if !checkIdentifier() {
        raiseError("invalid identifier")
        return false
    }

    if !checkOpeningBracket() { return false }

    if !compileParameterList() {
        raiseError("compileParameterList")
        return false
    }

    if !checkClosingBracket() { return false }

    writeOpen("subroutineBody")

    if !checkOpeningBrace() { return false }

    //possible variable declaration
    for checkToken("var") {
        if !compileVarDec() {
            raiseError("compileVarDec")
            return false
        }
    }

    //compile statements
    if !compileStatements() {
        raiseError("compileStatements")
        return false
    }

    if !checkClosingBrace() { return false }

    writeClose("subroutineBody")
    writeClose("subroutineDec")

    return true
}

func compileParameterList() bool {
    writeOpen("parameterList")

    //no parameters
    if checkToken(")") {
        writeClose("parameterList")
        return true
    }

    if !checkTypeAndIdentifier() {
        raiseError("checkTypeAndIdentifier")
        return false
    }

    for checkToken(",") {
        writeToken()
        if !checkTypeAndIdentifier() {
            raiseError("checkTypeAndIdentifier")
            return false
        }
    }

    writeClose("parameterList")
    return true
}

func compileVarDec() bool {
    writeOpen("varDec")
    writeToken()

    if !checkTypeAndIdentifier() {
        raiseError("checkTypeAndIdentifier")
        return false
    }

    for checkToken(",") {
        writeToken()
        if !checkIdentifier() { return false }
    }

    if !checkSemicolon() { return false }

    writeClose("varDec")
    return true
}

func compileStatements() bool {
    writeOpen("statements")
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
    writeClose("statements")
    return true
}

func compileLet() bool {
    writeOpen("compileLet")
    writeToken()

    if !checkIdentifier() { return false }

    //potential expression
    if checkToken("[") {
        writeToken()

        if !compileExpression() {
            raiseError("compileExpression")
            return false
        }

        if !checkToken("]") {
            raiseError("missing closing ]")
            return false
        }
        writeToken()
    }

    //equals expression
    if !checkToken("=") {
        raiseError("missing =")
        return false
    }

    writeToken()

    if !compileExpression() {
        raiseError("compileExpression")
        return false
    }

    if !checkSemicolon() { return false }

    writeClose("letStatement")
    return true
}

func compileIf() bool {
    writeOpen("ifStatement")
    writeToken()

    if !checkOpeningBracket() { return false }
    if !compileExpression() {
        raiseError("compileExpression")
        return false
    }
    if !checkClosingBracket() { return false }
    if !checkOpeningBrace() { return false }
    if !compileStatements() {
        raiseError("compileStatements")
        return false
    }
    if !checkClosingBrace() { return false }
    if checkToken("else") {
        if !checkOpeningBrace() { return false }
        if !compileStatements() { return false }
        if !checkClosingBrace() { return false }
    }
    writeClose("ifStatement")
    return true
}

func compileWhile() bool {
    writeOpen("whileStatement")
    writeToken()

    if !checkOpeningBracket() { return false }
    if !compileExpression() {
        raiseError("compileExpression")
        return false
    }
    if !checkClosingBracket() { return false }
    if !checkOpeningBrace() { return false }
    if !compileStatements() {
        raiseError("compileStatements")
        return false }
    if !checkClosingBrace() { return false }

    writeClose("whileStatement")
    return true
}

func compileDo() bool {
    writeOpen("doStatement")
    writeToken()

    if !checkSubroutineCall() {
        raiseError("subroutineCall")
        return false
    }
    if !checkSemicolon() { return false }

    writeClose("doStatement")
    return true
}

func compileReturn() bool {
    writeOpen("returnStatement")
    writeToken()

    //expression is optional
    if checkToken(";") {
        writeToken()
        writeClose("returnStatement")
        return true
    }

    if !compileExpression() {
        raiseError("compileExpression")
        return true
    }
    if !checkSemicolon() { return false }

    writeClose("returnStatement")
    return true
}

func compileExpression() bool {
    writeOpen("expression")
    //an expression must contain at least one term
    if !compileTerm() {
        raiseError("compileTerm")
        return false
    }

    //(op term)* 
    op := true
    for op {
        for i := range operators {
            if strings.Compare(operators[i], current[1]) == 0 {
                op = true
                writeToken()
                if !compileTerm() {
                    raiseError("compileTerm")
                    return false
                }
                break
            }
        op = false
        }
    }

    writeClose("expression")
    return true
}

func compileTerm() bool {
    writeOpen("term")

    //int, string or keyword
    if checkTokenTypeSlice([]string{"integerConstant", "stringConstant", "keywordConstant"}) {
        writeToken()
        return true
    }

    //(expression)
    if checkToken("(") {
        writeToken()
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
            writeToken()
            if !compileTerm() {
                raiseError("compileTerm")
                return false
            }
            return true
        }
    }

    //can now only have varName | varName [expression] | subroutineCall
    //Al of these terms begin with varName => look ahead one token to
    //differentiate

    if checkNextToken("[") {
        if !checkIdentifier() { return false }
        writeToken()

        if !compileExpression() { 
            raiseError("compileExpression")
            return false
        }

        if !checkToken("]") {
            raiseError("missing ]")
            return false
        }
        writeToken()

        return true
    }

    if checkNextToken("(") || checkNextToken(".") {
        if !checkSubroutineCall() {
            raiseError("checkSubroutineCall")
            return false
        }
        return true
    }

    if !checkIdentifier() { return false }

    writeClose("term")
    return true
}

func checkSubroutineCall() bool {
    if !checkIdentifier() { return false }

    //could be subroutineName(blargh) | class/var.subroutine(blargh)

    if checkToken(".") {
        writeToken()
        if !checkIdentifier() { return false }
    }

    if !checkOpeningBracket() { return false }

    if !compileExpressionList() { 
        raiseError("compileExpressionList")
        return false
    }
    if !checkClosingBracket() { return false }

    return true
}

func compileExpressionList() bool {
    writeOpen("expressionList")
    //first off, check if we have an empty list
    if checkToken(")") {
        writeClose("expressionList")
        return true
    }

    //now we know we have to compile at least one expression
    if !compileExpression() {
        raiseError("compileExpression")
        return false
    }

    for !checkToken(")") {
        if !checkToken(",") {
            raiseError("missing comma")
            return false
        }
        writeToken()
        if !compileExpression() {
            raiseError("compileExpression")
            return false
        }
    }
    writeClose("expressionList")
    return true
}

func main () {

    //arguments
    args := os.Args
    filepath := args[1]
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

    //first routine to be called must be compileClass
    if !compileClass() {
        raiseError("unable to compile class")
        debug()
    }

    //write output to file
    f, err := os.Create(target)
    if err != nil {
        os.Exit(1)
    }
    defer f.Close()

    for i := range output {
        f.WriteString(output[i])
        f.WriteString("\n")
    }
}
