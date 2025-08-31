package main

import (
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"golang.org/x/term"
)

type LifeGameLoop struct {
}

func (lh *LifeGameLoop) Start(parameters *UsageParameters) {
	termbox.Init()

	exitMessage := ""

	defer func() {
		termbox.Close()
		if exitMessage != "" {
			println(exitMessage)
		}
	}()

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	universe := newUniverse(width, height, *parameters.population)

	keyCh := make(chan termbox.Event, 1)

	go func() {
		for {
			ev := termbox.PollEvent() // blocking
			keyCh <- ev
		}
	}()

	// Main loop (board redraw)
	tick := time.NewTicker(*parameters.sleep)
	defer tick.Stop()
	i := 0
	terminate := false
	noAliveIteration := 0
	for {
		select {
		case ev := <-keyCh:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				return // exit
			}
		case <-tick.C:
			i = i + 1
			result := universe.printTillResizeComplete(i + 1)
			if result == NoAlive {
				noAliveIteration++
			}

			if noAliveIteration == 2 {
				exitMessage = "No one left alive. Leaving the game."
				terminate = true
			} else {
				universe.nextStep()
			}
		}

		if i == *parameters.gens || terminate {
			break
		}
	}
}
