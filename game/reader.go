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
	"bufio"
	"log"
	"os"
)

func readFile(source *string) [][]bool {

	file, err := os.Open(*source)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	var matrix [][]bool
	maxCols := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '!' {
			continue
		}
		cols := len(line)
		if cols > maxCols {
			maxCols = cols
		}
		row := make([]bool, cols)
		for i, ch := range line {
			if ch == 'O' {
				row[i] = true
			} else if ch == '.' {
				row[i] = false
			}
		}
		matrix = append(matrix, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	rows := len(matrix)
	transposed := make([][]bool, maxCols)
	for i := range transposed {
		transposed[i] = make([]bool, rows)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < len(matrix[r]); c++ {
			transposed[c][r] = matrix[r][c]
		}
	}

	return transposed
}
