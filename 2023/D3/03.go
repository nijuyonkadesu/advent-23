package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const DEBUG = true 
const PART = 2 // Change to 2 to solve part 2

type EngineBufferLayout struct {
    oneRaw string
    twoRaw string
    threeRaw string 

    oneDigits [][]int 
    twoDigits [][]int 
    threeDigits [][]int 

    oneSpecials [][]int
    twoSpecials [][]int 
    threeSpecials [][]int 

    oneValids []int
    twoValids []int
    threeValids []int
}

type EngineScanner interface {
    alloc()
    shiftUp()
    insert(scramble string) 

    // PART 1
    validate() // TODO change name
    pieceTogether()

    // PART 2
    findGearRatio()
    gearsTogether()

    finalize()
    customPrint()
}

type EngineBuffer struct {
    buffer EngineBufferLayout
    scramble string // this one should be the older one
    serial int

    gearValue int
    gearRatio int
    addMulSwitch bool // false is add, true is mul
}

func (eb *EngineBuffer) alloc(len int) {
    eb.buffer.oneDigits = make([][]int, len)
    eb.buffer.twoDigits = make([][]int, len)
    eb.buffer.threeDigits = make([][]int, len)

    eb.buffer.oneSpecials = make([][]int, len)
    eb.buffer.twoSpecials= make([][]int, len)
    eb.buffer.threeSpecials = make([][]int, len)

    eb.buffer.oneValids = make([]int, len)
    eb.buffer.twoValids = make([]int, len)
    eb.buffer.threeValids = make([]int, len)

    eb.gearRatio = 0
    eb.addMulSwitch = true 
}

func (eb *EngineBuffer) customPrint() {

    fmt.Println("D:", eb.buffer.oneDigits)
    fmt.Println("D:", eb.buffer.twoDigits)
    fmt.Println("D:", eb.buffer.threeDigits)

    fmt.Println(eb.buffer.oneRaw)
    fmt.Println(eb.buffer.twoRaw)
    fmt.Println(eb.buffer.threeRaw)

    fmt.Println("S:", eb.buffer.oneSpecials)
    fmt.Println("S:", eb.buffer.twoSpecials)
    fmt.Println("S:", eb.buffer.threeSpecials)

    fmt.Println("V:", eb.buffer.oneValids)
    fmt.Println("V:", eb.buffer.twoValids)
    fmt.Println("V:", eb.buffer.threeValids)

    fmt.Println("------------------------------", eb.serial)
}

func (eb *EngineBuffer) shiftUp() {
    eb.buffer.oneRaw, eb.buffer.twoRaw = eb.buffer.twoRaw, eb.buffer.threeRaw
    eb.buffer.oneDigits, eb.buffer.twoDigits = eb.buffer.twoDigits, eb.buffer.threeDigits
    eb.buffer.oneSpecials, eb.buffer.twoSpecials = eb.buffer.twoSpecials, eb.buffer.threeSpecials
    eb.buffer.oneValids, eb.buffer.twoValids = eb.buffer.twoValids, eb.buffer.threeValids
}

func (eb *EngineBuffer) pieceTogether() {
    for i, v := range eb.buffer.oneDigits {
        if eb.buffer.oneValids[i] == 1 {
            numStr := eb.buffer.oneRaw[v[0]:v[1]]
            num, err := strconv.Atoi(numStr)
            if err != nil {
                fmt.Println("pieceTogether:", err)
            }
            if DEBUG {
                fmt.Println("+", num)
            }
            eb.serial = eb.serial + num
        }
    }
}

func (eb *EngineBuffer) validate() {
    for _, v := range eb.buffer.twoSpecials {
        if len(v) == 0 {
            break;
        }
        splIdx := v[0]

        for i, numIdx := range eb.buffer.oneDigits {
            if splIdx >= numIdx[0] && splIdx < numIdx[1] {
                eb.buffer.oneValids[i] = 1
                
            } else if splIdx + 1 < len(eb.buffer.oneRaw) && splIdx + 1 >= numIdx[0] && splIdx + 1 < numIdx[1] {
                eb.buffer.oneValids[i] = 1
                
            } else if splIdx - 1 > 0 && splIdx - 1 >= numIdx[0] && splIdx - 1 < numIdx[1] {
                eb.buffer.oneValids[i] = 1
                
            }
        }

        for i, numIdx := range eb.buffer.twoDigits {
            if splIdx >= numIdx[0] && splIdx < numIdx[1] {
                eb.buffer.twoValids[i] = 1
            } else if splIdx + 1 < len(eb.buffer.twoRaw) && splIdx + 1 >= numIdx[0] && splIdx + 1 < numIdx[1] {
                eb.buffer.twoValids[i] = 1
                
            } else if splIdx - 1 > 0 && splIdx - 1 >= numIdx[0] && splIdx - 1 < numIdx[1] {
                eb.buffer.twoValids[i] = 1
                
            }
        }

        for i, numIdx := range eb.buffer.threeDigits {
            if splIdx >= numIdx[0] && splIdx < numIdx[1] {
                eb.buffer.threeValids[i] = 1
            } else if splIdx + 1 < len(eb.buffer.threeRaw) && splIdx + 1 >= numIdx[0] && splIdx + 1 < numIdx[1] {
                eb.buffer.threeValids[i] = 1
                
            } else if splIdx - 1 > 0 && splIdx - 1 >= numIdx[0] && splIdx - 1 < numIdx[1] {
                eb.buffer.threeValids[i] = 1
                
            }
        }
    }
    eb.pieceTogether()
}

func (eb *EngineBuffer) gearsTogether() {
    for i, v := range eb.buffer.oneDigits {
        if eb.buffer.oneValids[i] == 1 {
            numStr := eb.buffer.oneRaw[v[0]:v[1]]
            num, err := strconv.Atoi(numStr)
            if err != nil {
                fmt.Println("gearsTogether:", err)
            }
            if DEBUG {
                fmt.Println("*", num)
            }
            eb.addMulSwitch = !eb.addMulSwitch // false at first run
            if eb.addMulSwitch {
                eb.gearRatio = eb.gearRatio + eb.gearValue * num
            } else {
                eb.gearValue = num 
            }
        }
    }
}

func (eb *EngineBuffer) findGearRatio() {
    for _, v := range eb.buffer.twoSpecials {
        if len(v) == 0 {
            break;
        }
        splIdx := v[0]
        matches := 0
        rowOneMatch := []int {} 
        rowTwoMatch := []int {} 
        rowThreeMatch := []int {}
        // TODO: one star with only 2 adjacent number. can be at any row
        // Change the below logic, it's wrong
        // Backtrack ? woah

        for i, numIdx := range eb.buffer.oneDigits {
            if splIdx >= numIdx[0] && splIdx < numIdx[1] {
                rowOneMatch = append(rowOneMatch, i)
                matches++
            } else if splIdx + 1 < len(eb.buffer.oneRaw) && splIdx + 1 >= numIdx[0] && splIdx + 1 < numIdx[1] {
                rowOneMatch = append(rowOneMatch, i)
                matches++
            } else if splIdx - 1 > 0 && splIdx - 1 >= numIdx[0] && splIdx - 1 < numIdx[1] {
                rowOneMatch = append(rowOneMatch, i)
                matches++
            }
        }

        for i, numIdx := range eb.buffer.twoDigits {
            if splIdx >= numIdx[0] && splIdx < numIdx[1] {
                rowTwoMatch = append(rowTwoMatch, i)
                matches++
            } else if splIdx + 1 < len(eb.buffer.twoRaw) && splIdx + 1 >= numIdx[0] && splIdx + 1 < numIdx[1] {
                rowTwoMatch = append(rowTwoMatch, i)
                matches++
            } else if splIdx - 1 > 0 && splIdx - 1 >= numIdx[0] && splIdx - 1 < numIdx[1] {
                rowTwoMatch = append(rowTwoMatch, i)
                matches++
            }
        }

        for i, numIdx := range eb.buffer.threeDigits {
            if splIdx >= numIdx[0] && splIdx < numIdx[1] {
                rowThreeMatch = append(rowThreeMatch, i)
                matches++
            } else if splIdx + 1 < len(eb.buffer.threeRaw) && splIdx + 1 >= numIdx[0] && splIdx + 1 < numIdx[1] {
                rowThreeMatch = append(rowThreeMatch, i)
                matches++
            } else if splIdx - 1 > 0 && splIdx - 1 >= numIdx[0] && splIdx - 1 < numIdx[1] {
                rowThreeMatch = append(rowThreeMatch, i)
                matches++
            }
        }

        if matches == 2 {
            for _, v := range rowOneMatch {
                if len(rowOneMatch) != 0 {
                    eb.buffer.oneValids[v] = 1
                }
            }

            for _, v := range rowTwoMatch {
                if len(rowTwoMatch) != 0 {
                    eb.buffer.twoValids[v] = 1
                }
            }

            for _, v := range rowThreeMatch {
                if len(rowThreeMatch) != 0 {
                    eb.buffer.threeValids[v] = 1
                }
            }
        } 
    }
    eb.gearsTogether()
}

func (eb *EngineBuffer) finalize(){
    eb.shiftUp()
    if PART == 1 {
        eb.validate()
    } else if PART == 2 {
        eb.findGearRatio()

    }
    eb.shiftUp()
    if PART == 1 {
        eb.validate()
    } else if PART == 2 {
        eb.findGearRatio()
    }
}

func regexPattern() (*regexp.Regexp, *regexp.Regexp) {

    if PART == 1 {
        return regexp.MustCompile(`[^\w.]`), regexp.MustCompile(`\d+`)
    } else if PART == 2 {
        return regexp.MustCompile(`\*`), regexp.MustCompile(`\d+`)
    }
        return regexp.MustCompile(`[^\w.]`), regexp.MustCompile(`\d+`)
}

func (eb *EngineBuffer) insert(scramble string) {
    eb.scramble = scramble

    reSpl, reNums := regexPattern()

    spl := reSpl.FindAllIndex([]byte(scramble), -1)
    nums := reNums.FindAllIndex([]byte(scramble), -1)

     // if DEBUG {
     //     fmt.Println(scramble)
     //     fmt.Println("RS:", spl)
     //     fmt.Println("RN:", nums)
     // }

    eb.shiftUp()
    eb.buffer.threeRaw = eb.scramble
    eb.buffer.threeDigits = nums
    eb.buffer.threeSpecials = spl
    eb.buffer.threeValids = make([]int, len(eb.buffer.threeDigits))

    if PART == 1 {
        eb.validate()
    } else if PART == 2 {
        eb.findGearRatio()
    }
}

func main () {
    if PART == 1 {
        fmt.Println(parseSerialParts())
    } else if PART == 2 {
        fmt.Println(parseGearParts())
    }
}

func parseSerialParts() int {
    engineBuffer := EngineBuffer{}

    file, err := os.Open("03.txt")
    if err != nil {
        fmt.Println(err)
        return 0
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for i := 0; scanner.Scan(); i++ {
        scramble := scanner.Text()
        if scramble == "" {
            continue
        }
        if i == 0 {
            engineBuffer.alloc(len(scramble))
        }
        engineBuffer.insert(scramble)

        if DEBUG {
            engineBuffer.customPrint()
        }
        
        if err != nil {
            fmt.Print(err)
        }
    }

    engineBuffer.finalize()
    return engineBuffer.serial
}

func parseGearParts() int {
    engineBuffer := EngineBuffer{}

    file, err := os.Open("tiny.txt")
    if err != nil {
        fmt.Println(err)
        return 0
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for i := 0; scanner.Scan(); i++ {
        scramble := scanner.Text()
        if scramble == "" {
            continue
        }
        if i == 0 {
            engineBuffer.alloc(len(scramble))
        }
        engineBuffer.insert(scramble)

        if DEBUG {
            engineBuffer.customPrint()
        }
        
        if err != nil {
            fmt.Print(err)
        }
    }

    engineBuffer.finalize()
    return engineBuffer.gearRatio
}
