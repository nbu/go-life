package main

import (
	//"bufio"

	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"golang.org/x/term"
)

type Universe struct {
	width     int
	height    int
	board     [][]bool
	nextBoard [][]bool
}

func main() {

	termbox.Init()
	defer termbox.Close()

	//wordScanner := bufio.NewScanner(os.Stdin)
	//wordScanner.Split(bufio.ScanWords)
	//wordScanner.Scan()
	//var height, _ = strconv.ParseInt(wordScanner.Text(), 10, 32)
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	universe := newUniverse(width, height)

	gens := 100

	for i := 0; i < gens; i++ {
		universe.printUniverse(i + 1)
		time.Sleep(100 * time.Millisecond)
		universe.makeStep()
	}
}

func newUniverse(width int, height int) *Universe {
	universe := new(Universe)
	universe.board = make([][]bool, width)
	universe.nextBoard = make([][]bool, width)
	universe.width = width
	universe.height = height

	for i := range universe.board {
		universe.board[i] = make([]bool, height)
		universe.nextBoard[i] = make([]bool, height)

		for j := range universe.board[i] {
			universe.board[i][j] = rand.Intn(10) == 1
			universe.nextBoard[i][j] = false
		}
	}

	return universe
}

func (u *Universe) printUniverse(gen int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	//	var s string
	var alive int
	// TODO Borys - this width, height safeguard from terminal resizing by the user is not working
	// instead check at the end if terminal was changed, if it was changed regenerate
	// term.GetSize itself is expensive and slow everything significantly if call on each iteration
	//width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	for i := range u.board {
		//if i > width {
		//	break
		//}
		for j := range u.board[i] {
			//if j > height {
			//	break
			//}
			var cell rune
			if u.board[i][j] {
				cell = 'O'
				alive++
			} else {
				cell = ' '
			}
			//			s += cell
			termbox.SetCell(i, j, cell, termbox.ColorGreen, termbox.ColorDefault)
		}
		//		s += "\n"
	}

	//	fmt.Print("\033[H\033[2J")
	//	fmt.Printf("Generation #%d\n", gen)
	//	fmt.Printf("Alive: %d\n", alive)
	//	fmt.Printf(s)
	termbox.Flush()
}

func (u *Universe) makeStep() {

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
