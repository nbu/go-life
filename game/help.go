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
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

const DefaultSymbolAlive = 'O'

type UsageParameters struct {
	gens        *int
	population  *int
	sleep       *time.Duration
	file        *string
	symbolAlive rune
	boardType   *string
}

type LifeHelp struct {
}

func (lh *LifeHelp) DefineUsage() *UsageParameters {
	usageParameters := new(UsageParameters)

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Game of Life Simulator\n\n")
		fmt.Fprintf(os.Stderr, "This program simulates Conway's Game of Life on a terminal grid.\n")
		fmt.Fprintf(os.Stderr, "You can control generations, population density, speed, initial layout file, board type (infinite or boarded).\n")
		fmt.Fprintf(os.Stderr, "In the ininite board mode you can pan the board with the arrow keys. Also you can use mouse wheel to scroll up and down. To reset origin back pres 'r'.\n\n")
		fmt.Fprintf(os.Stderr, "To pause simulation press <SPACE>.\n\n")
		fmt.Fprintf(os.Stderr, "To end simulation at any time press <ESC>.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		pflag.PrintDefaults()
	}

	usageParameters.gens =
		pflag.IntP(
			"gens",
			"g",
			1000,
			"number of generations to run simulation, 0 means infinite number of generations")
	usageParameters.population =
		pflag.IntP(
			"population",
			"p",
			20,
			"population density\nif initial layout file is provided this parameter is ignored\n")
	usageParameters.sleep =
		pflag.DurationP(
			"sleep",
			"s", 100*time.Millisecond,
			"number of seconds to sleep between generations")
	usageParameters.file =
		pflag.StringP(
			"file",
			"f",
			"",
			"initial layout file")
	symbolAlive :=
		pflag.StringP("symbol-alive",
			"a",
			string(DefaultSymbolAlive),
			"symbol to represent alive cell on the board\nunicode character can be provided as $'\\u2591'")
	usageParameters.boardType =
		pflag.StringP("board-type",
			"t",
			"infinite",
			"board type to simulate, allowed values are infinite or boarded")
	pflag.Parse()

	if len(*symbolAlive) > 0 {
		usageParameters.symbolAlive = []rune(*symbolAlive)[0]
	} else {
		usageParameters.symbolAlive = DefaultSymbolAlive
	}

	return usageParameters
}
