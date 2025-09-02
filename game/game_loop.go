package game

import (
	"time"

	"github.com/nsf/termbox-go"
)

const speedIncrement = 5 * time.Millisecond

type LifeGameLoop struct {
}

func (lh *LifeGameLoop) Start(parameters *UsageParameters) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	exitMessage := ""

	defer func() {
		termbox.Close()
		if exitMessage != "" {
			println(exitMessage)
		}
	}()

	universe := NewUniverse(parameters)

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
	terminate := false
	resetTimer := false
	for {
		select {
		case ev := <-keyCh:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					return
				} else if ev.Key == termbox.KeyArrowLeft {
					universe.Pan(-1, 0)
				} else if ev.Key == termbox.KeyArrowRight {
					universe.Pan(1, 0)
				} else if ev.Key == termbox.KeyArrowUp {
					universe.Pan(0, -1)
				} else if ev.Key == termbox.KeyArrowDown {
					universe.Pan(0, 1)
				} else if ev.Ch == '-' {
					*parameters.sleep = *parameters.sleep + speedIncrement
					resetTimer = true
				} else if ev.Ch == '+' {
					*parameters.sleep = *parameters.sleep - speedIncrement
					resetTimer = true
					if *parameters.sleep < speedIncrement {
						*parameters.sleep = speedIncrement
					}
				} else if ev.Ch == 'r' {
					universe.ResetOrigin(Coord{0, 0})
				}
			}
		case <-tick.C:
			PrintTillResizeComplete(universe)

			if universe.AliveCount() == 0 {
				exitMessage = "Extinction of the population"
				terminate = true
			} else {
				universe.NextStep()
			}
		}

		if (universe.Generation() == *parameters.gens && *parameters.gens > 0) || terminate {
			break
		}

		if resetTimer {
			tick.Reset(*parameters.sleep)
			resetTimer = false
		}
	}
}
