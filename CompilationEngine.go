package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "strings"
    "regexp"
)

var tokens [][]string

func generateTokenArray(slice []string) {

    for i := range slice {
        blargh := strings.Split(slice[i], " ")
        tokens = append(tokens, blargh)
    }
    fmt.Println(tokens[0][0])
    fmt.Println(tokens[0][1])
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
}
