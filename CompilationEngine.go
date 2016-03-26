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

func compileClass() bool {

    if checkToken("class") {
        writeOpen("class")
        writeToken()
    }
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
    if !(compileClass()) {
        raiseError("unable to compile class")
        return
    }
}
