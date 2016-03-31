package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "strings"
    "regexp"
)

var tokens [][]string
var current []string
var output []string

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

func compileClass() bool {

    if !checkToken("class") {
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

    return true
}

func compileClassVarDec() bool {
    writeOpen("classVarDec")
    writeToken()

    //check type
    if (!checkTokenSlice([]string{"int, char, bool"}) && !checkIdentifier()) {
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

func main () {
    //open tokens file
    data, err := ioutil.ReadFile("tokens")
    if err != nil {
        os.Exit(1)
    }
    stringified := string(data)
    re, _ := regexp.Compile("\"")
    stringified = re.ReplaceAllString(stringified, "")
    slice := strings.Split(stringified, "\n")
    generateTokenArray(slice)
    //first routine to be called must be compileClass
    if !compileClass() {
        raiseError("unable to compile class")
        return
    }
    fmt.Println(output)
}
