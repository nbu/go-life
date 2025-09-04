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

	game := NewGame(parameters)

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
					game.Pan(-1, 0)
				} else if ev.Key == termbox.KeyArrowRight {
					game.Pan(1, 0)
				} else if ev.Key == termbox.KeyArrowUp {
					game.Pan(0, -1)
				} else if ev.Key == termbox.KeyArrowDown {
					game.Pan(0, 1)
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
					game.ResetOrigin(Coord{0, 0})
				} else if ev.Key == termbox.KeySpace {
					pause = !pause
				}
			}
		case <-tick.C:
			game.PrintTillResizeComplete()

			if pause {
				continue
			}

			if game.Universe.AliveCount() == 0 {
				exitMessage = "Extinction of the population"
				terminate = true
			} else {
				game.Universe.NextStep()
			}
		}

		if (game.Universe.Generation() == *parameters.gens && *parameters.gens > 0) || terminate {
			break
		}

		if resetTimer {
			tick.Reset(*parameters.sleep)
			resetTimer = false
		}
	}
}
