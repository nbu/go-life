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
