package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type UsageParameters struct {
	gens       *int
	population *int
	sleep      *time.Duration
}

type LifeHelp struct {
}

func (lh *LifeHelp) DefineUsage() *UsageParameters {
	usageParameters := new(UsageParameters)

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Game of Life Simulator\n\n")
		fmt.Fprintf(os.Stderr, "This program simulates Conway's Game of Life on a terminal grid.\n")
		fmt.Fprintf(os.Stderr, "You can control generations, population density, and speed.\n")
		fmt.Fprintf(os.Stderr, "To end simulation at any time press <ESC>.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		pflag.PrintDefaults()
	}

	usageParameters.gens =
		pflag.IntP(
			"gens",
			"g",
			1000,
			"number of generations")
	usageParameters.population =
		pflag.IntP(
			"population",
			"p",
			20,
			"percentage of population (whole board filled is 100%)")
	usageParameters.sleep =
		pflag.DurationP(
			"sleep",
			"s", 10*time.Millisecond,
			"number of seconds to sleep between generations")

	pflag.Parse()

	return usageParameters
}
