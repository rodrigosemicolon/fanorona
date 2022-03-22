package main

import (
	"fmt"

	"github.com/rodrigosemicolon/fanorona/game"
)

func main() {
	testGame := game.NewBoard()
	var x, y, newX, newY int
	isOver := 0
	for isOver == 0 {
		testGame.PrintBoard()
		fmt.Println()
		fmt.Println("captures available")
		caps := testGame.CapturesAvailable()
		for cap := range caps {

			fmt.Println("from: ", caps[cap].InitialPos, "to: ", caps[cap].EndingPos)
			fmt.Println("captures: ", caps[cap].CappedPieces)

		}

		fmt.Println("Player: ", testGame.Turn)
		fmt.Println("select piece row")
		fmt.Scan(&x)
		fmt.Println("select piece col")
		fmt.Scan(&y)

		fmt.Println("select piece row")
		fmt.Scan(&newX)
		fmt.Println("select piece col")
		fmt.Scan(&newY)
		if x == -1 && testGame.LastMove.Player == testGame.Turn {
			testGame.Turn = -testGame.Turn

		} else {

			mv := testGame.NewMove(
				game.Pos{X: x, Y: y},
				game.Pos{X: newX, Y: newY},
			)
			if mv.Invalid == nil {
				isOver = testGame.ApplyMove(mv)

			} else {
				fmt.Println(mv.Invalid)
			}
		}

	}

}
