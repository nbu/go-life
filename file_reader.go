package main

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
	defer file.Close()

	var matrix [][]bool
	maxCols := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cols := len(line)
		if cols > maxCols {
			maxCols = cols
		}
		row := make([]bool, cols)
		for i, ch := range line {
			if ch == 'x' {
				row[i] = true
			} else {
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
