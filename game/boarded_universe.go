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
}

func CreateUniverseBoarded(width int, height int, parameters *UsageParameters) *BoardedUniverse {

	u := new(BoardedUniverse)
	u.board = make([][]int, width)
	u.nextBoard = make([][]int, width)
	u.width = width
	u.height = height
	u.parameters = *parameters

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

	if x >= 0 && x < u.width && y >= 0 && y < u.height {
		u.board[x][y] = u.board[x][y] + 1
		u.aliveCount++
	}
}

func (u *BoardedUniverse) NextStep() {

	aliveCount := 0
	for i := range u.board {
		for j := range u.board[i] {
			u.nextBoard[i][j] = u.isAliveOnNextStep(i, j)
			if u.nextBoard[i][j] > 0 {
				aliveCount++
			}
		}
	}

	tmpBoard := u.board
	u.board = u.nextBoard
	u.nextBoard = tmpBoard
	u.aliveCount = aliveCount
	u.generation++
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

func (u *BoardedUniverse) IsAlive(i int, j int) int {

	return u.isAlive(i, j, false)
}

func (u *BoardedUniverse) Pan(_ int, _ int) {
}

func (u *BoardedUniverse) AliveCount() int {
	return u.aliveCount
}

func (u *BoardedUniverse) Generation() int {
	return u.generation
}
