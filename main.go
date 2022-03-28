package main

import (
	"bufio"
	"os"

	"github.com/rodrigosemicolon/fanorona/game"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	testGame := game.NewGame(reader)
	game.RunGame(testGame)
}
