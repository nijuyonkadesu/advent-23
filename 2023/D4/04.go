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

const DEBUG = false 
const PART = 1 

type Card struct {
    winningNumbers []int 
    myNumbers []int 
    matches int
}

type CardUtils interface {
    init()
    computeMatches()
    totalJackpot()
}

type InvestmentPortfolio struct {
    luckyCard Card // dunno what to name, it's the last card 'me' looked up
    thisBatchCardReturns float64
}

type InvestmentPortfolioUtils interface {
    noteLuckyCardReward()
    clearRewards() // Not implemented yet
}

func (c *Card) init(win []string, nums []string)  {
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

func main()  {
    if PART == 1 {
        fmt.Println(examinePile())
    }
}

// DONE: parse and save as 2 arrays (winningNumbers, myNumbers)
// DONE: Check with winningNumbers - find just the count
// DONE: maintain a score start with 0 =  2^0 (1), 2^1 (2), 2^2 (4), 2^3 (8) -> 2 power (count-1)
// DONE: Sum everything

func examinePile() float64 {
    portfolio := InvestmentPortfolio {}
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

        card.init(winningNumbers, myNumbers)
        card.computeMatches()
        portfolio.luckyCard = card
        portfolio.noteLuckyCardReward()
    }

    return portfolio.thisBatchCardReturns
}


