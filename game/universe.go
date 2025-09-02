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
type Coord [2]int
type Bounds [2]Coord

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
	GameBounds() Bounds
	Origin() Coord
	ResetOrigin(coord Coord)
}

func NewUniverse(parameters *UsageParameters) Universe {

	var u Universe
	screenWidth, screenHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	screenWidth -= 2
	screenHeight -= 2
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
		termbox.SetCell(i, 0, '\u2500', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(i, height-1, '\u2500', termbox.ColorDefault, termbox.ColorDefault)
	}

	for i := range height {
		termbox.SetCell(0, i, '\u2502', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(width-1, i, '\u2502', termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.SetCell(0, 0, '\u250C', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(width-1, 0, '\u2510', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(0, height-1, '\u2514', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(width-1, height-1, '\u2518', termbox.ColorDefault, termbox.ColorDefault)

	for i := range width - 2 {
		for j := range height - 2 {
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
			termbox.SetCell(i+1, j+1, cell, fgColor, termbox.ColorDefault)
		}
	}

	bounds := u.GameBounds()
	origin := u.Origin()
	if bounds[0][0] < origin[0] {
		termbox.SetCell(0, height/2, '\u25C0', termbox.ColorDefault, termbox.ColorDefault)
	}
	if bounds[1][0] > origin[0]+width-3 {
		termbox.SetCell(width-1, height/2, '\u25B6', termbox.ColorDefault, termbox.ColorDefault)
	}
	if bounds[0][1] < origin[1] {
		termbox.SetCell(width/2, 0, '\u25B2', termbox.ColorDefault, termbox.ColorDefault)
	}
	if bounds[1][1] > origin[1]+height-3 {
		termbox.SetCell(width/2, height-1, '\u25BC', termbox.ColorDefault, termbox.ColorDefault)
	}

	generationsText := fmt.Sprintf(" Generation: %d ", u.Generation())
	drawString(
		2,
		height-1,
		generationsText,
		termbox.ColorDefault,
		termbox.ColorDefault)

	originText := fmt.Sprintf(" Origin: x=%d y=%d", u.Origin()[0], u.Origin()[1])
	drawString(
		2,
		0,
		originText,
		termbox.ColorDefault,
		termbox.ColorDefault)
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

func drawString(x, y int, s string, fg, bg termbox.Attribute) {
	for i, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}
