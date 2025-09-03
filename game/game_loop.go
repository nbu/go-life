package game

import (
	"time"

	"github.com/nsf/termbox-go"
)

const SpeedIncrement = 5 * time.Millisecond

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
	pause := false
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
					*parameters.sleep = *parameters.sleep + SpeedIncrement
					resetTimer = true
				} else if ev.Ch == '+' {
					*parameters.sleep = *parameters.sleep - SpeedIncrement
					resetTimer = true
					if *parameters.sleep < SpeedIncrement {
						*parameters.sleep = SpeedIncrement
					}
				} else if ev.Ch == 'r' {
					universe.ResetOrigin(Coord{0, 0})
				} else if ev.Key == termbox.KeySpace {
					pause = !pause
				}
			}
		case <-tick.C:
			PrintTillResizeComplete(universe)

			if pause {
				continue
			}

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
