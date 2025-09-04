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
	AliveCount() int
	Generation() int
	GameBounds() Bounds
	Stats() map[int]UniverseStats
}

type Game struct {
	Universe Universe
	Origin   Coord
}

func NewGame(parameters *UsageParameters) Game {

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

	game := Game{
		Universe: u,
		Origin:   Coord{0, 0},
	}

	if *parameters.file != "" {
		matrix := readFile(parameters.file)
		game.embedMatrix(matrix, screenWidth, screenHeight)
	}

	return game
}

func (game *Game) embedMatrix(source [][]bool, screenWidth int, screenHeight int) {

	if len(source) > screenWidth || len(source[0]) > screenHeight && *game.Universe.Parameters().boardType == "boarded" {
		log.Fatal("Source matrix is larger than target matrix")
		return
	}

	colOffset := (screenWidth - len(source)) / 2
	rowOffset := (screenHeight - len(source[0])) / 2

	for r, row := range source {
		for c, val := range row {
			if val {
				game.Universe.SetAliveCell(colOffset+r, rowOffset+c)
			}
		}
	}
}

func (game *Game) Pan(x int, y int) {
	game.Origin.X = game.Origin.X + x
	game.Origin.Y = game.Origin.Y + y
}

func (game *Game) ResetOrigin(coord Coord) {
	game.Origin = coord
}

func (game *Game) PrintTillResizeComplete() {
	for {
		result := game.printUniverse()
		if result != BoardResized {
			return
		}
	}
}

func (game *Game) printUniverse() BoardPrintResult {
	err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		panic(err)
	}

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	game.drawBorder(width, height)
	game.drawCells(width, height)
	game.drawNavigationArrows(height, width)
	game.drawInfoText(height, width)

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

func (game *Game) drawInfoText(height int, width int) {
	u := game.Universe

	stats := u.Stats()
	genStats := stats[u.Generation()]

	generationsText := fmt.Sprintf(" Generation: %d; Population: %d ",
		u.Generation(),
		genStats.alive)
	drawString(
		2,
		height-1,
		generationsText,
		termbox.ColorDefault,
		termbox.ColorDefault)

	originText := fmt.Sprintf(" Origin: x=%d y=%d", game.Origin.X, game.Origin.Y)
	drawString(
		2,
		0,
		originText,
		termbox.ColorDefault,
		termbox.ColorDefault)

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

func (game *Game) drawCells(width int, height int) {
	u := game.Universe

	for i := range width - 2 {
		for j := range height - 2 {
			var cell rune
			isAlive := u.IsAlive(i+game.Origin.X, j+game.Origin.Y)
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

func (game *Game) drawBorder(width int, height int) {
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

func (game *Game) drawNavigationArrows(height int, width int) {

	u := game.Universe

	bounds := u.GameBounds()
	origin := game.Origin
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
