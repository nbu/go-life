package main

import (
	"life/framework"
	"life/game"
)

func main() {

	help := new(game.LifeHelp)
	gameLoop := new(game.LifeGameLoop)
	framework.Run(help, gameLoop)
}
