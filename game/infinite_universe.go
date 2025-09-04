/*
 * Copyright (c) 2025 Borys Nebosenko
 *
 * This file is part of Go-life.
 *
 * Go-life is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published
 * by the Free Software Foundation, either version 3 of the License,
 * or (at your option) any later version.
 *
 * Go-life is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Go-life.  If not, see <https://www.gnu.org/licenses/>.
 */

package game

import (
	"math"
)

type InfiniteUniverse struct {
	Universe
	board      map[Coord]int
	parameters UsageParameters
	generation int
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
			neighbor := Coord{c.X + n.X, c.Y + n.Y}
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
	width := u.bounds.BottomRight.X - u.bounds.TopLeft.X
	height := u.bounds.BottomRight.Y - u.bounds.TopLeft.Y
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
	return u.board[Coord{x, y}]
}

func (u *InfiniteUniverse) Parameters() UsageParameters {
	return u.parameters
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

func (u *InfiniteUniverse) Stats() map[int]UniverseStats {
	return u.stats
}

func (u *InfiniteUniverse) setBounds(cell Coord) {
	if cell.X < u.bounds.TopLeft.X {
		u.bounds.TopLeft.X = cell.X
	}

	if cell.Y < u.bounds.TopLeft.Y {
		u.bounds.TopLeft.Y = cell.Y
	}

	if cell.X > u.bounds.BottomRight.X {
		u.bounds.BottomRight.X = cell.X
	}

	if cell.Y > u.bounds.BottomRight.Y {
		u.bounds.BottomRight.Y = cell.Y
	}
}

func (u *InfiniteUniverse) resetBounds() {
	u.bounds = Bounds{
		Coord{math.MaxInt, math.MaxInt},
		Coord{math.MinInt, math.MinInt},
	}
}
