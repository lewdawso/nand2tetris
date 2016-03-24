package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "strings"
    "regexp"
)

var tokens [][]string
var output []string

func generateTokenArray(slice []string) {

    for i := range slice {
        pair := strings.Split(slice[i], " ")
        tokens = append(tokens, pair)
    }
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

func raiseError(message string) {
    fmt.Println("error: ", message)
}

func checkToken(slice []string) bool {
    fmt.Println(slice)
}

func compileClass() bool {

    //if checkToken("class") {
    //    writeOpen("class")
    //}
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
    writeOpen("class")
    checkToken("blargh")
}
