package game

import (
	"math"
)

type InfiniteUniverse struct {
	Universe
	board      map[Coord]int
	parameters UsageParameters
	generation int
	origin     Coord
	bounds     Bounds
	stats      map[int]UniverseStats
}

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
	u.resetBounds()
	u.stats = make(map[int]UniverseStats)

	return u
}

func (u *InfiniteUniverse) NextStep() {
	newBoard := make(map[Coord]int)
	counts := make(map[Coord]int)
	u.generation++

	u.resetBounds()
	for c := range u.board {
		for _, n := range neighbors {
			neighbor := Coord{c[0] + n[0], c[1] + n[1]}
			counts[neighbor]++
		}
	}

	stats := UniverseStats{}
	for cell, cnt := range counts {
		if cnt == 3 || (cnt == 2 && u.board[cell] > 0) {
			if u.board[cell] == 0 {
				stats.born++
			}
			newBoard[cell] = u.board[cell] + 1
			u.setBounds(cell)
		} else {
			if u.board[cell] > 0 {
				stats.died++
			}
		}
	}

	stats.alive = u.AliveCount()
	u.setDeadCount(stats)
	u.stats[u.generation] = stats

	u.board = newBoard
}

func (u *InfiniteUniverse) setDeadCount(stats UniverseStats) {
	width := u.bounds[1][0] - u.bounds[0][0]
	height := u.bounds[1][1] - u.bounds[0][1]
	stats.dead += width*height - stats.alive
}

func (u *InfiniteUniverse) SetAliveCell(x int, y int) {
	oldStatus := u.IsAlive(x, y)
	coord := Coord{x, y}
	u.board[coord] = u.board[coord] + 1
	u.setBounds(coord)
	u.setStats(oldStatus, 1)
}

func (u *InfiniteUniverse) setStats(oldStatus int, aliveInc int) {

	stats := u.stats[u.generation]
	stats.alive = stats.alive + aliveInc
	u.setDeadCount(stats)
	stats.died += oldStatus
	u.stats[u.generation] = stats
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

func (u *InfiniteUniverse) GameBounds() Bounds {
	return u.bounds
}

func (u *InfiniteUniverse) Origin() Coord {
	return u.origin
}

func (u *InfiniteUniverse) ResetOrigin(coord Coord) {
	u.origin = coord
}

func (u *InfiniteUniverse) Stats() map[int]UniverseStats {
	return u.stats
}

func (u *InfiniteUniverse) setBounds(cell Coord) {
	if cell[0] < u.bounds[0][0] {
		u.bounds[0][0] = cell[0]
	}

	if cell[1] < u.bounds[0][1] {
		u.bounds[0][1] = cell[1]
	}

	if cell[0] > u.bounds[1][0] {
		u.bounds[1][0] = cell[0]
	}

	if cell[1] > u.bounds[1][1] {
		u.bounds[1][1] = cell[1]
	}
}

func (u *InfiniteUniverse) resetBounds() {
	u.bounds = Bounds{
		Coord{math.MaxInt, math.MaxInt},
		Coord{math.MinInt, math.MinInt},
	}
}
