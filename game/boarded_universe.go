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

type BoardedUniverse struct {
	Universe
	width      int
	height     int
	board      [][]int
	nextBoard  [][]int
	parameters UsageParameters
	aliveCount int
	generation int
	bounds     Bounds
	stats      map[int]UniverseStats
}

func CreateUniverseBoarded(width int, height int, parameters *UsageParameters) *BoardedUniverse {

	u := new(BoardedUniverse)
	u.board = make([][]int, width)
	u.nextBoard = make([][]int, width)
	u.width = width
	u.height = height
	u.parameters = *parameters
	u.bounds = Bounds{Coord{0, 0}, Coord{width - 1, height - 1}}
	u.stats = make(map[int]UniverseStats)

	for i := range u.board {
		u.board[i] = make([]int, height)
		u.nextBoard[i] = make([]int, height)
	}

	return u
}

func (u *BoardedUniverse) Parameters() UsageParameters {
	return u.parameters
}

func (u *BoardedUniverse) SetAliveCell(x int, y int) {

	oldStatus := u.board[x][y]
	if x >= 0 && x < u.width && y >= 0 && y < u.height {
		u.board[x][y] = u.board[x][y] + 1
		u.aliveCount++
	}
	u.setStats(oldStatus, 1)
}

func (u *BoardedUniverse) setStats(oldStatus int, aliveInc int) {

	stats := u.stats[u.generation]
	stats.alive = stats.alive + aliveInc
	u.setDeadCount(stats)
	stats.died += oldStatus
	u.stats[u.generation] = stats
}

func (u *BoardedUniverse) setDeadCount(stats UniverseStats) {
	width := u.bounds.BottomRight.X - u.bounds.TopLeft.X
	height := u.bounds.BottomRight.Y - u.bounds.TopLeft.Y
	stats.dead += width*height - stats.alive
}

func (u *BoardedUniverse) NextStep() {

	stats := UniverseStats{}
	aliveCount := 0
	u.generation++
	for i := range u.board {
		for j := range u.board[i] {
			isAlive := u.isAliveOnNextStep(i, j)
			if isAlive > 0 && u.board[i][j] == 0 {
				stats.born++
			} else if isAlive == 0 && u.board[i][j] > 0 {
				stats.died++
			}
			u.nextBoard[i][j] = isAlive
			if u.nextBoard[i][j] > 0 {
				aliveCount++
			}
		}
	}

	stats.alive = u.AliveCount()
	u.setDeadCount(stats)
	u.stats[u.generation] = stats

	tmpBoard := u.board
	u.board = u.nextBoard
	u.nextBoard = tmpBoard
	u.aliveCount = aliveCount
}

func (u *BoardedUniverse) isAliveOnNextStep(i int, j int) int {

	cnt := u.aliveNeighbours(i, j)
	if cnt == 3 || (cnt == 2 && u.board[i][j] > 0) {
		return u.board[i][j] + 1
	} else {
		return 0
	}
}

func (u *BoardedUniverse) aliveNeighbours(x int, y int) int {
	cnt := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (i != 0 || j != 0) && u.isAlive(x+i, y+j, true) > 0 {
				cnt++
			}
		}
	}

	return cnt
}

func (u *BoardedUniverse) isAlive(i int, j int, wrapEdges bool) int {

	if !wrapEdges && (i < 0 || i >= u.width || j < 0 || j >= u.height) {
		return 0
	}

	i += u.width
	j += u.height

	return u.board[i%u.width][j%u.height]
}

func (u *BoardedUniverse) IsAlive(x int, y int) int {

	return u.isAlive(x, y, false)
}

func (u *BoardedUniverse) AliveCount() int {
	return u.aliveCount
}

func (u *BoardedUniverse) Generation() int {
	return u.generation
}

func (u *BoardedUniverse) GameBounds() Bounds {
	return u.bounds
}

func (u *BoardedUniverse) Stats() map[int]UniverseStats {
	return u.stats
}
