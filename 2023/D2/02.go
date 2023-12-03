package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const DEBUG = false 
const RED_MAX = 12
const GREEN_MAX = 13
const BLUE_MAX = 14

type Game struct {
    red []int
    green []int
    blue []int
    id int 
}

func main() {
    fmt.Println(dontLieJoshua())
}

func (g *Game) isPossible() bool {
    for _, v := range g.red {
        if v > RED_MAX {
            return false 
        }
    }
    for _, v := range g.green {
        if v > GREEN_MAX {
            return false 
        }
    }
    for _, v := range g.blue {
        if v > BLUE_MAX {
            return false 
        }
    }

    return true
}

func (g *Game) addCubes(count int, color string) {
    switch color {
    case "red":
        g.red = append(g.red, count)
    case "green":
        g.green = append(g.green, count)
    case "blue":
        g.blue = append(g.blue, count)
    }
}

func dontLieJoshua() int {
    sumOfGameIds := 0
    file, err := os.Open("02.txt")
    if err != nil {
        fmt.Println(err)
        return 0
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        claim := scanner.Text()
        if claim == "" {
            continue
        }
        game, err := toGame(claim)
        if err != nil {
            fmt.Print(err)
        }
        if game.isPossible() {
            sumOfGameIds = sumOfGameIds + game.id
        }
    }
    return sumOfGameIds
}

func toGame(claim string) (Game, error){
    game := Game {}

    re := regexp.MustCompile(`\b\w+`)
    res := re.FindAllString(claim, -1)
    id, err := strconv.Atoi(res[1])
    if err != nil {
        return game, err 
    }
    game.id = id 

    if DEBUG {
        fmt.Println(res)
    }

    for i := 2; i < len(res); i+=2 {
        count, err := strconv.Atoi(res[i]) 
        if err != nil {
            return game, err
        }
        color := res[i+1]
        game.addCubes(count, color)
    }

    if DEBUG {
        fmt.Println(game)
    }
    return game, nil
}

