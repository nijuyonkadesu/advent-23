package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)
const DEBUG = false 

func main()  {

    if DEBUG {
        test()
    } else {
        parseDocument()
    }
}

func test() {

    fmt.Print(findCallibrationNaive("ok8ok9"))
    fmt.Print(findCallibrationPro("oneight"))
}

func parseDocument() {
    sum := 0
    file, err := os.Open("01.txt")
    if err != nil {
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        doc := scanner.Text()
        sum += findCallibrationPro(doc)
    }
    fmt.Print(sum)
}

// Part 1
func findCallibrationNaive(source string) int {

    re := regexp.MustCompile(`\d`)
    res := re.FindAllString(source, -1)
    if res == nil {
        return 0
    }
    if DEBUG {
        fmt.Println(res)
    }
    tenth, _ := strconv.Atoi(res[0])
    ones, _ := strconv.Atoi(res[len(res)-1])
    return tenth*10 + ones
}

// Part 2
func findCallibrationPro(source string) int {

    translator := map[string] int {
        "one": 1, 
        "two": 2, 
        "three": 3, 
        "four": 4,
        "five": 5, 
        "six": 6, 
        "seven": 7,
        "eight": 8,
        "nine": 9, 
        "zero": 0,
        "1": 1, 
        "2": 2, 
        "3": 3, 
        "4": 4,
        "5": 5, 
        "6": 6, 
        "7": 7,
        "8": 8,
        "9": 9, 
        "0": 0,
    }

    forwardPass := regexp.MustCompile(`\d|one|two|three|four|five|six|seven|eight|nine|zero`)
    backwaredPass := regexp.MustCompile(`orez|enin|thgie|neves|xis|evif|ruof|eerht|owt|eno|\d`)
    tenth := forwardPass.FindAllString(source, 1)
    if tenth == nil {
        return 0
    }
    ones := reverseString(backwaredPass.FindAllString(reverseString(source), 1)[0])
    res := translator[tenth[0]]*10 + translator[ones]
    if DEBUG {
        fmt.Println(tenth[0] + ones)
    }
    return res
}

func reverseString(s string) string {

    runes := []rune(s)
    length := len(runes)
    for s, e := 0, length-1; s < e; s, e = s+1, e-1 {
        runes[s], runes[e] = runes[e], runes[s]
    }
    return string(runes)
}
