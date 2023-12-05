package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
)

const DEBUG = true 
const PART = 2 

type Card struct {
    id int
    winningNumbers []int 
    myNumbers []int 
    matches int
}

type CardUtils interface {
    init()
    computeMatches()
    totalJackpot()
}

type Dupes struct {
    flag bool
    count int 
}

type InvestmentPortfolio struct {
    luckyCard Card // dunno what to name, it's the last card 'me' looked up
    thisBatchCardReturns float64
    totalCards int
    cardMultiplier int
    thisBatchInstanceValidity map[int]Dupes
}

type InvestmentPortfolioUtils interface {
    noteLuckyCardReward()

    updateMultiplier()
    updateCardCount()
}

func (c *Card) init(win, nums []string, id int)  {
    c.id = id
    c.winningNumbers = make([]int, len(win))
    c.myNumbers = make([]int, len(nums))
    c.matches = 0

    for i, v := range win {
        num, err  := strconv.Atoi(v)
        if err != nil {
            fmt.Println(err)
        }
        c.winningNumbers[i] = num 
    }

    for i, v := range nums {
        num, err  := strconv.Atoi(v)
        if err != nil {
            fmt.Println(err)
        }
        c.myNumbers[i] = num 
    }
}

func (c *Card) computeMatches()  {
    for _, v := range c.myNumbers {
        idx := slices.Index(c.winningNumbers, v)
        if idx != -1 {
            c.matches++
            if DEBUG {
                fmt.Print(".",v)
            }
        }
    }
    if DEBUG {
        fmt.Println(" match:", c.matches)
    }
}

func (c *Card) totalJackpot() float64 {
    power := c.matches - 1
    if power == -1 {
        return 0.0
    }
    jackpot := math.Pow(2.0, float64(power)) 
    if DEBUG {
        fmt.Println("jackpot:", jackpot)
    }
    return jackpot
}

func (i *InvestmentPortfolio) noteLuckyCardReward() {
    i.thisBatchCardReturns += i.luckyCard.totalJackpot()
}

func (i *InvestmentPortfolio) updateMultiplier()  {
    i.updateCardCount()
    if i.luckyCard.matches > 0 { 
        end := i.luckyCard.id + i.luckyCard.matches
        dupes := i.thisBatchInstanceValidity[end]
        dupes.flag = true
        dupes.count++
        i.thisBatchInstanceValidity[end] = dupes
        if DEBUG { fmt.Println("new till", dupes) }

        i.cardMultiplier++
        if DEBUG { fmt.Println("ADD", dupes.count, "id:", i.luckyCard.id,  "new mul:", i.cardMultiplier, "till:", end) }
    }
    if i.thisBatchInstanceValidity[i.luckyCard.id].flag { 
        dupes := i.thisBatchInstanceValidity[i.luckyCard.id]
        i.cardMultiplier -= dupes.count 
        if DEBUG { fmt.Println("SUB", dupes.count, "id:", i.luckyCard.id) }
    }
    if DEBUG { fmt.Println("final mul:", i.cardMultiplier) }
    if DEBUG { fmt.Println("-------------------") }
}

func (i *InvestmentPortfolio) updateCardCount()  {
    i.totalCards += i.cardMultiplier * i.luckyCard.matches
    if DEBUG { fmt.Println("tot", i.totalCards) }
}

func main()  {
    if PART == 1 {
        fmt.Println(examinePile())
    } else if PART == 2 {
        fmt.Println(examinePileCardCount())
    }
}

// DONE: parse and save as 2 arrays (winningNumbers, myNumbers)
// DONE: Check with winningNumbers - find just the count
// DONE: maintain a score start with 0 =  2^0 (1), 2^1 (2), 2^2 (4), 2^3 (8) -> 2 power (count-1)
// DONE: Sum everything

// DESC : 4 matching at card 1 -> you win 2, 3, 4, 5 -> one copy each (coz first occurance)
// 2 instance of card 2 -> 2 matches -> you win 3, 4 -> two copy each
// 4 instace of card 3 ............................ -> four copy each
// compute original + copies count 

// DONE: Track card number
// DONE: multiplicative factor = 1, and increase when match is found
// DONE: have a ~HEAP~ map to track end times
// DONE: multiplicative -1 when a match is found on heap (dupes is allowed)
// TODO: call update fns at appropriate places


func examinePile() float64 {
    portfolio := InvestmentPortfolio { cardMultiplier: 1 }
    card := Card{}

    file, err := os.Open("04.txt")
    if err != nil {
        fmt.Println(err)
        return 0.0
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        colorfulCard := scanner.Text()
        if len(colorfulCard) == 0 {
            continue;
        }

        re := regexp.MustCompile(`\d+`)
        cardNums := re.FindAllString(colorfulCard, -1)

        // 1:6, 6: for tiny.txt
        winningNumbers := cardNums[1:11]
        myNumbers := cardNums[11:]
        if DEBUG {
            fmt.Println("w:", winningNumbers)
            fmt.Println("m:", myNumbers)
        }

        card.init(winningNumbers, myNumbers, 0)
        card.computeMatches()
        portfolio.luckyCard = card
        portfolio.noteLuckyCardReward()
    }

    return portfolio.thisBatchCardReturns
}

func examinePileCardCount() int {
    portfolio := InvestmentPortfolio { 
        cardMultiplier: 1,
        thisBatchInstanceValidity: make(map[int]Dupes),
    }
    card := Card{}
    cardsSeen := 0

    file, err := os.Open("tiny.txt")
    if err != nil {
        fmt.Println(err)
        return 0
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        colorfulCard := scanner.Text()
        if len(colorfulCard) == 0 {
            continue;
        }

        re := regexp.MustCompile(`\d+`)
        cardNums := re.FindAllString(colorfulCard, -1)

        // 1:6, 6: for tiny.txt
        winningNumbers := cardNums[1:6]
        myNumbers := cardNums[6:]
        if DEBUG {
            fmt.Println("w:", winningNumbers)
            fmt.Println("m:", myNumbers)
        }

        cardsSeen++
        card.init(winningNumbers, myNumbers, cardsSeen)
        card.computeMatches()
        portfolio.luckyCard = card
        // portfolio.updateCardCount()
        portfolio.updateMultiplier()
    }

    return portfolio.totalCards
}
