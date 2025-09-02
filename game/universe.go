package game

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/nsf/termbox-go"
	"golang.org/x/term"
)

type BoardPrintResult int

const (
	Printed BoardPrintResult = iota
	BoardResized
)

type Universe interface {
	SetAliveCell(x int, y int)
	IsAlive(x int, y int) int
	Parameters() UsageParameters
	NextStep()
	// Pan TODO Borys: Pan function doesn't belong to the game, it's a representation not a game logic
	Pan(x int, y int)
	AliveCount() int
	Generation() int
}

func NewUniverse(screenWidth int, screenHeight int, parameters *UsageParameters) Universe {

	var u Universe
	if *parameters.boardType == "infinite" {
		u = CreateUniverseInfinite(parameters)
	} else if *parameters.boardType == "boarded" {
		u = CreateUniverseBoarded(screenWidth, screenHeight, parameters)
	} else {
		fmt.Printf("Invalid board-type specified: %s\n", *parameters.boardType)
		os.Exit(3)
	}

	for i := range screenWidth {

		if *parameters.file != "" {
			continue
		}

		for j := range screenHeight {
			if rand.Intn(100) <= *parameters.population {
				u.SetAliveCell(i, j)
			}
		}
	}

	if *parameters.file != "" {
		matrix := readFile(parameters.file)
		embedMatrix(matrix, screenWidth, screenHeight, u)
	}

	return u
}

func embedMatrix(source [][]bool, screenWidth int, screenHeight int, u Universe) {

	// Check source fits
	if len(source) > screenWidth || len(source[0]) > screenHeight {
		log.Fatal("Source matrix is larger than target matrix")
		return
	}

	colOffset := (screenWidth - len(source)) / 2
	rowOffset := (screenHeight - len(source[0])) / 2

	// Copy source into target
	for r, row := range source {
		for c, val := range row {
			if val {
				u.SetAliveCell(colOffset+r, rowOffset+c)
			}
		}
	}
}

func PrintTillResizeComplete(u Universe) {
	for {
		result := printUniverse(u)
		if result != BoardResized {
			return
		}
	}
}

func printUniverse(u Universe) BoardPrintResult {
	err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		panic(err)
	}
	var alive int
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	for i := range width {
		for j := range height {
			var cell rune
			isAlive := u.IsAlive(i, j)
			if isAlive > 0 {
				cell = u.Parameters().symbolAlive
				alive++
			} else {
				cell = ' '
			}

			var fgColor termbox.Attribute

			if isAlive > 1 {
				fgColor = termbox.ColorDarkGray
			} else {
				fgColor = termbox.ColorGreen
			}
			termbox.SetCell(i, j, cell, fgColor, termbox.ColorDefault)
		}
	}
	newWidth, newHeight, _ := term.GetSize(int(os.Stdout.Fd()))

	result := Printed
	if newWidth != width || newHeight != height {
		result = BoardResized
	} else {
		result = Printed
	}

	if result != BoardResized {
		err := termbox.Flush()
		if err != nil {
			panic(err)
		}
	}

	return result
}
