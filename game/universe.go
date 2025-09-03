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
type Coord struct {
	X int
	Y int
}
type Bounds struct {
	TopLeft     Coord
	BottomRight Coord
}

const (
	Printed BoardPrintResult = iota
	BoardResized
)

type UniverseStats struct {
	alive int
	born  int
	dead  int
	died  int
}

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
	Stats() map[int]UniverseStats
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
	if len(source) > screenWidth || len(source[0]) > screenHeight && *u.Parameters().boardType == "boarded" {
		log.Fatal("Source matrix is larger than target matrix")
		return
	}

	colOffset := (screenWidth - len(source)) / 2
	rowOffset := (screenHeight - len(source[0])) / 2

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

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	drawBorder(width, height)
	drawCells(u, width, height)
	drawNavigationArrows(u, height, width)
	drawInfoText(u, height, width)

	result := Printed

	newWidth, newHeight, _ := term.GetSize(int(os.Stdout.Fd()))
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

func drawInfoText(u Universe, height int, width int) {
	generationsText := fmt.Sprintf(" Generation: %d ", u.Generation())
	drawString(
		2,
		height-1,
		generationsText,
		termbox.ColorDefault,
		termbox.ColorDefault)

	originText := fmt.Sprintf(" Origin: x=%d y=%d", u.Origin().X, u.Origin().Y)
	drawString(
		2,
		0,
		originText,
		termbox.ColorDefault,
		termbox.ColorDefault)

	stats := u.Stats()
	genStats := stats[u.Generation()]
	trend := 0.0
	if genStats.died > 0 {
		trend = float64(genStats.born) / float64(genStats.died)
	}

	statsText := fmt.Sprintf(" Born: %d; Died: %d; Born/Died: %f ",
		genStats.born, genStats.died, trend)
	drawString(
		width-2-len(statsText),
		height-1,
		statsText,
		termbox.ColorDefault,
		termbox.ColorDefault)

	bounds := u.GameBounds()
	sizeText := fmt.Sprintf(" Size: width=%d height=%d ",
		bounds.BottomRight.X-bounds.TopLeft.X, bounds.BottomRight.Y-bounds.TopLeft.Y)
	drawString(
		width-2-len(sizeText),
		0,
		sizeText,
		termbox.ColorDefault,
		termbox.ColorDefault)
}

func drawCells(u Universe, width int, height int) {
	for i := range width - 2 {
		for j := range height - 2 {
			var cell rune
			isAlive := u.IsAlive(i, j)
			if isAlive > 0 {
				cell = u.Parameters().symbolAlive
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
}

func drawBorder(width int, height int) {
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
}

func drawNavigationArrows(u Universe, height int, width int) {

	bounds := u.GameBounds()
	origin := u.Origin()
	if bounds.TopLeft.X < origin.X {
		termbox.SetCell(0, height/2, '\u25C0', termbox.ColorDefault, termbox.ColorDefault)
	}
	if bounds.BottomRight.X > origin.X+width-3 {
		termbox.SetCell(width-1, height/2, '\u25B6', termbox.ColorDefault, termbox.ColorDefault)
	}
	if bounds.TopLeft.Y < origin.Y {
		termbox.SetCell(width/2, 0, '\u25B2', termbox.ColorDefault, termbox.ColorDefault)
	}
	if bounds.BottomRight.Y > origin.Y+height-3 {
		termbox.SetCell(width/2, height-1, '\u25BC', termbox.ColorDefault, termbox.ColorDefault)
	}
}

func drawString(x, y int, s string, fg, bg termbox.Attribute) {
	for i, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}
