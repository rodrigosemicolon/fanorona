package main

import (
	"fmt"

	"github.com/rodrigosemicolon/fanorona/game"
)

func main() {
	testGame := game.NewBoard()
	var x, y, newX, newY int
	var canSkip bool
	player := 1
	for {
		testGame.PrintBoard()
		fmt.Println()
		fmt.Println("Player: ", player)
		fmt.Println("select piece row,col and insert destination row,col:")
		fmt.Scanf("%d,%d %d,%d", &x, &y, &newX, &newY)
		fmt.Println(x)
		if x == -1 && canSkip {
			player = -player
			canSkip = false
			continue
		}
		mv := testGame.NewMove(player,
			game.Pos{X: x, Y: y},
			game.Pos{X: newX, Y: newY},
		)
		if mv.Invalid == nil {
			player, canSkip = testGame.ApplyMove(mv)

		} else {
			fmt.Println(mv.Invalid)
		}

	}

}
