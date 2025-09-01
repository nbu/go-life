package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

const defaultSymbolAlive = 'O'

type UsageParameters struct {
	gens         *int
	population   *int
	sleep        *time.Duration
	file         *string
	symbol_alive rune
}

type LifeHelp struct {
}

func (lh *LifeHelp) DefineUsage() *UsageParameters {
	usageParameters := new(UsageParameters)

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Game of Life Simulator\n\n")
		fmt.Fprintf(os.Stderr, "This program simulates Conway's Game of Life on a terminal grid.\n")
		fmt.Fprintf(os.Stderr, "You can control generations, population density, speed and initial layout file.\n")
		fmt.Fprintf(os.Stderr, "To end simulation at any time press <ESC>.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		pflag.PrintDefaults()
	}

	usageParameters.gens =
		pflag.IntP(
			"gens",
			"g",
			1000,
			"number of generations, 0 means infinite number of generations")
	usageParameters.population =
		pflag.IntP(
			"population",
			"p",
			20,
			"percentage of the random population\nif initial layout file is provided this parameter is ignored\n")
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
	symbol_alive :=
		pflag.StringP("symbol-alive",
			"a",
			string(defaultSymbolAlive),
			"symbol to represent alive cell on the board\nunicode character must be provided as $'\\u2591'")
	pflag.Parse()

	if len(*symbol_alive) > 0 {
		usageParameters.symbol_alive = []rune(*symbol_alive)[0]
	} else {
		usageParameters.symbol_alive = defaultSymbolAlive
	}

	return usageParameters
}
