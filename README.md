# Go Game of Life

A fast and scalable implementation of Conway’s Game of Life in Go running in the terminal.  
Supports both **bounded boards** and **sparse infinite boards**.

## Features

- Infinite board mode (sparse storage, scales to very large grids).
- Bounded board mode with fixed width/height determined by the widht and height of the terminal.
- Terminal rendering (via [termbox-go](https://github.com/nsf/termbox-go)).
- Keyboard controls for pausing and adjusting speed.
  - pause: \<SPACE\>
  - quit: \<ESC\>
  - speed up/down: +/-
  - pan board with arrows: left, right, up and down 
  - reset board origin: r
- Unicode characters for smooth board visualization.
- Statistics tracking:
    - Generation count
    - Alive cells
    - Born / Died cells per tick

## Demo
<video src="doc/go-life.mp4" controls>
  Your browser does not support the video tag.
</video>