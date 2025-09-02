package game

type Coord [2]int

type InfiniteUniverse struct {
	Universe
	board      map[Coord]int
	parameters UsageParameters
	origin     Coord
	generation int
}

// Neighbors around a cell
var neighbors = [8]Coord{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func CreateUniverseInfinite(parameters *UsageParameters) *InfiniteUniverse {
	u := new(InfiniteUniverse)
	u.board = make(map[Coord]int)
	u.parameters = *parameters
	u.origin = Coord{0, 0}

	return u
}

func (u *InfiniteUniverse) NextStep() {
	newBoard := make(map[Coord]int)
	counts := make(map[Coord]int)

	// Count neighbors for each live cell
	for c := range u.board {
		for _, n := range neighbors {
			neighbor := Coord{c[0] + n[0], c[1] + n[1]}
			counts[neighbor]++
		}
	}

	// Apply rules
	for cell, cnt := range counts {
		if cnt == 3 || (cnt == 2 && u.board[cell] > 0) {
			newBoard[cell] = u.board[cell] + 1
		}
	}

	u.board = newBoard
	u.generation++
}

func (u *InfiniteUniverse) SetAliveCell(x int, y int) {
	u.board[Coord{x, y}] = u.board[Coord{x, y}] + 1
}

func (u *InfiniteUniverse) IsAlive(x int, y int) int {
	x = x + u.origin[0]
	y = y + u.origin[1]
	return u.board[Coord{x, y}]
}

func (u *InfiniteUniverse) Parameters() UsageParameters {
	return u.parameters
}

func (u *InfiniteUniverse) Pan(x int, y int) {
	u.origin[0] = u.origin[0] + x
	u.origin[1] = u.origin[1] + y
}

func (u *InfiniteUniverse) AliveCount() int {
	return len(u.board)
}

func (u *InfiniteUniverse) Generation() int {
	return u.generation
}
