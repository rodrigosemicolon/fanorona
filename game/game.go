package game

import (
	"bufio"
	"errors"
	"fmt"
	"math"
)

type Game struct {
	Board       *Board
	History     []Move
	Turn        int
	inputReader *bufio.Reader
}

func NewGame(input *bufio.Reader) *Game {
	g := Game{
		Board:       NewBoard(),
		History:     []Move{},
		Turn:        1,
		inputReader: input,
	}
	return &g
}

func (game *Game) CapturesAvailable() []Move {
	captures := make([]Move, 0)
	for i, row := range game.Board.BoardState {
		for j, piece := range row {
			if piece == game.Turn {
				for hOffset := -1; hOffset <= 1; hOffset++ {
					for vOffset := -1; vOffset <= 1; vOffset++ {
						m := game.NewMove(Pos{i, j}, Pos{i + hOffset, j + vOffset})
						if m.Invalid == nil {
							if m.IsCapture {
								captures = append(captures, m)
							}
						}
					}
				}
			}
		}
	}
	return captures
}

func (game *Game) NewRecapture(first Pos, second Pos) Move {
	var recap Move
	previousMove := game.History[len(game.History)-1]
	if first.X == previousMove.EndingPos.X && first.Y == previousMove.EndingPos.Y {
		recap = game.NewMove(first, second)
	} else {
		recap = Move{Invalid: errors.New("Must recapture with the same piece")}
	}
	return recap
}

func (game *Game) PromptMove() (Pos, Pos) {
	fmt.Println(game.Board.ToString())

	fmt.Println("insert move as: startingX,startingY endingX,endingY")
	input, err := game.inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	var startX, startY, endX, endY int
	fmt.Sscanf(input, "%d,%d %d,%d", &startX, &startY, &endX, &endY)
	return Pos{startX, startY}, Pos{endX, endY}

}

func (game *Game) PromptRecapture() *Pos {
	fmt.Println(game.Board.ToString())

	fmt.Println("insert move as: nextX,nextY")
	fmt.Println("-1,-1 will skip turn")
	input, err := game.inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	var nextX, nextY int
	fmt.Sscanf(input, "%d,%d", &nextX, &nextY)
	if nextX == -1 || nextY == -1 {
		return nil
	}
	return &Pos{nextX, nextY}

}

func (game *Game) PromptDirection() int {
	fmt.Println("to capture forward insert 1")
	fmt.Println("to capture backward insert 0")

	input, err := game.inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	var decision int
	fmt.Sscanf(input, "%d", &decision)
	return decision
}

func RunGame(game *Game) {
	var moves [2]Pos
	var history = make([]Move, 0)
	var isOver int

	for i := 0; isOver == 0; i++ {
		moves[0], moves[1] = game.PromptMove()

		Mv := game.NewMove(moves[0], moves[1])
		if Mv.Invalid == nil {
			isOver = game.ApplyMove(Mv)
			history = append(history, Mv)
			for Mv.IsCapture {
				if len(Mv.CappedFwd) > 0 && len(Mv.CappedBwd) > 0 {
					choice := game.PromptDirection()
					if choice == 1 { //true means capture forward
						Mv.CappedBwd = nil
					} else { // false means capture backward
						Mv.CappedFwd = nil

					}
				}
				moves[0] = moves[1]
				newPos := game.PromptRecapture()

				//if equals(Mv.EndingPos, moves) {
				if newPos == nil {
					break
				}
				moves[1] = *newPos
				newCap := game.NewRecapture(moves[0], moves[1])
				if newCap.Invalid != nil {
					Mv.ReCaptures = append(Mv.ReCaptures, newCap)
					isOver = game.ApplyMove(newCap)

				}
			}
		}

	}

}

//create a new move if the positions are valid
func (game *Game) NewMove(mFrom Pos, mTo Pos) Move {

	if !game.Board.IsValid(mFrom) {
		return Move{Invalid: errors.New("Must select a piece on the board")}
	}

	if !game.Board.IsValid(mTo) {
		return Move{Invalid: errors.New("Pieces can't move out of bounds")}
	}

	if game.Turn != game.Board.BoardState[mFrom.X][mFrom.Y] {
		return Move{Invalid: errors.New("Players must move their own pieces")}

	}

	if game.Board.BoardState[mTo.X][mTo.Y] != 0 {
		return Move{Invalid: errors.New("Can't move a piece to a place already containing another")}
	}

	horizontalShift := math.Abs(float64(mFrom.X - mTo.X))
	verticalShift := math.Abs(float64(mFrom.Y - mTo.Y))

	if horizontalShift < 0 || horizontalShift > 1 {
		return Move{Invalid: errors.New("Can only move one place at a time")}
	}
	if verticalShift < 0 || verticalShift > 1 {
		return Move{Invalid: errors.New("Can only move one place at a time")}
	}
	if (mFrom.X+mFrom.Y)%2 == 0 {
		//strong intersection
		if 1 <= (horizontalShift+verticalShift) && (horizontalShift+verticalShift) <= 2 {
			m := Move{
				Player:     game.Turn,
				InitialPos: mFrom,
				EndingPos:  mTo,
				Invalid:    nil,
			}
			// TODO: decide capture front or back
			game.Board.CheckCaptures(&m)
			return m
		} else {
			return Move{Invalid: errors.New("Invalid move")}
		}
	} else {
		//weak intersection
		if horizontalShift+verticalShift == 1 {
			m := Move{
				Player:     game.Turn,
				InitialPos: mFrom,
				EndingPos:  mTo,
				Invalid:    nil,
			}
			game.Board.CheckCaptures(&m)
			return m

		} else {
			return Move{Invalid: errors.New("Invalid move")}
		}
	}

}

//applying a move
//returning an int representing whether the game has ended and who won
func (game *Game) ApplyMove(m Move) int {
	game.Board.BoardState[m.InitialPos.X][m.InitialPos.Y] = 0
	game.Board.BoardState[m.EndingPos.X][m.EndingPos.Y] = m.Player
	game.History = append(game.History, m)
	//if no captures its the other players turn
	if !m.IsCapture {
		game.Turn = -m.Player

	} else {
		if m.CappedFwd == nil {
			for _, piece := range m.CappedBwd {
				game.Board.BoardState[piece.X][piece.Y] = 0
			}

		} else {
			for _, piece := range m.CappedFwd {
				game.Board.BoardState[piece.X][piece.Y] = 0
			}
		}
		//CHECK WIN
		white, black := game.Board.CheckPieces()
		if white == 0 || black == 0 {
			return m.Player
		}
		//MISSING CHECK IF NO LEGAL MOVES AVAILABLE

	}
	return 0

}
