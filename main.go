package main

import (
	"life/framework"
)

func main() {

	help := new(LifeHelp)
	gameLoop := new(LifeGameLoop)
	framework.Run(help, gameLoop)
}
