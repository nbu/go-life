package main

import (
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
	NoAlive
)

type Universe struct {
	width      int
	height     int
	board      [][]bool
	nextBoard  [][]bool
	parameters UsageParameters
}

func newUniverse(width int, height int, parameters *UsageParameters) *Universe {
	universe := new(Universe)
	universe.board = make([][]bool, width)
	universe.nextBoard = make([][]bool, width)
	universe.width = width
	universe.height = height
	universe.parameters = *parameters

	for i := range universe.board {
		universe.board[i] = make([]bool, height)
		universe.nextBoard[i] = make([]bool, height)

		if *parameters.file != "" {
			continue
		}

		for j := range universe.board[i] {
			universe.board[i][j] = rand.Intn(100) <= *parameters.population
			universe.nextBoard[i][j] = false
		}
	}

	if *parameters.file != "" {
		matrix := readFile(parameters.file)
		embedMatrix(matrix, universe.board)
	}

	return universe
}

func embedMatrix(source [][]bool, target [][]bool) {

	targetRows, targetCols := len(target), len(target[0])

	// Check source fits
	if len(source) > targetRows || len(source[0]) > targetCols {
		log.Fatal("Source matrix is larger than target matrix")
		return
	}

	// Compute offsets to center source
	rowOffset := (targetRows - len(source)) / 2
	colOffset := (targetCols - len(source[0])) / 2

	// Copy source into target
	for r, row := range source {
		for c, val := range row {
			target[rowOffset+r][colOffset+c] = val
		}
	}
}

func (u *Universe) printTillResizeComplete(gen int) BoardPrintResult {
	for {

		result := u.printUniverse(gen)
		if result != BoardResized {
			return result
		}
	}
}

func (u *Universe) printUniverse(gen int) BoardPrintResult {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	var alive int
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	for i := range u.board {
		for j := range u.board[i] {
			var cell rune
			if u.board[i][j] {
				cell = u.parameters.symbol_alive
				alive++
			} else {
				cell = ' '
			}
			termbox.SetCell(i, j, cell, termbox.ColorGreen, termbox.ColorDefault)
		}
	}
	newWidth, newHeight, _ := term.GetSize(int(os.Stdout.Fd()))

	result := Printed
	if newWidth != width || newHeight != height {
		result = BoardResized
	} else if alive == 0 {
		result = NoAlive
	} else {
		result = Printed
	}

	if result != BoardResized {
		termbox.Flush()
	}

	return result
}

func (u *Universe) nextStep() {

	for i := range u.board {
		for j := range u.board[i] {
			u.nextBoard[i][j] = u.isAliveOnNextStep(i, j)
		}
	}

	tmpBoard := u.board
	u.board = u.nextBoard
	u.nextBoard = tmpBoard
}

func (u *Universe) isAliveOnNextStep(i int, j int) bool {

	cnt := u.aliveNeighbours(i, j)
	if !u.board[i][j] {
		return cnt == 3
	} else {
		return (cnt == 2) || (cnt == 3)
	}
}

func (u *Universe) aliveNeighbours(x int, y int) int {
	cnt := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (i != 0 || j != 0) && u.isAlive(x+i, y+j) {
				cnt++
			}
		}
	}

	return cnt
}

func (u *Universe) isAlive(i int, j int) bool {

	i += u.width
	j += u.height

	return u.board[i%u.width][j%u.height]
}
