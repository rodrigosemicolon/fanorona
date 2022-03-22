package main

import (
	"fmt"

	"github.com/rodrigosemicolon/fanorona/game"
)

func main() {
	testGame := game.NewBoard()
	var x, y, newX, newY int
	inputChan := make(chan [2]game.Pos)
	outputChan := make(chan string)
	testGame.PrintBoard()
	go func() {
		game.Game(testGame, inputChan, outputChan)
	}()
	for {
		/*
			caps := testGame.CapturesAvailable()
			for cap := range caps {

				fmt.Println("from: ", caps[cap].InitialPos, "to: ", caps[cap].EndingPos)
				fmt.Println("captures: ", caps[cap].CappedPieces)

			}
		*/

		fmt.Println("select piece row")
		fmt.Scan(&x)
		if x == -1 {
			inputChan <- [2]game.Pos{{X: -1, Y: -1}, {X: -1, Y: -1}}

		} else {
			fmt.Println("select piece col")
			fmt.Scan(&y)

			fmt.Println("select piece row")
			fmt.Scan(&newX)
			fmt.Println("select piece col")
			fmt.Scan(&newY)

			inputChan <- [2]game.Pos{{X: x, Y: y}, {X: newX, Y: newY}}
			fmt.Println(<-outputChan)
		}

	}

}
