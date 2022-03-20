package game

import (
	"errors"
	"fmt"
	"math"
)

//general board game state struct
type Board struct {
	BoardState [5][9]int
	Rows       int
	Cols       int
}

//move struct
type Move struct {
	Player       int
	InitialPos   Pos
	EndingPos    Pos
	IsCapture    bool
	CappedPieces []Pos
	Invalid      error
}

//board position struct
type Pos struct {
	X int
	Y int
}

//create a default starting board
func NewBoard() *Board {
	board := Board{
		BoardState: [5][9]int{
			{-1, -1, -1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1, -1, -1},
			{-1, 1, -1, 1, 0, 1, -1, 1, -1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		Rows: 5,
		Cols: 9,
	}
	return &board

}

//create a Board instance from a given position
func NewBoardFromArray(state [5][9]int) *Board {
	board := Board{
		BoardState: state,
		Rows:       5,
		Cols:       9,
	}
	return &board
}

//check the amount of pieces for both players
func (board *Board) CheckPieces() (int, int) {
	var whitePieces int
	var blackPieces int
	for _, row := range board.BoardState {
		for _, piece := range row {
			if piece == 1 {
				whitePieces++
			} else if piece == -1 {
				blackPieces++
			}
		}
	}
	return whitePieces, blackPieces
}

//check if a position falls within the borders of the game board
func (board *Board) IsValid(position Pos) bool {
	if 0 <= position.X && position.X < board.Rows && 0 <= position.Y && position.Y < board.Cols {
		return true
	} else {
		return false
	}
}

//check if any pieces are captured following a certain move
func (board *Board) CheckCaptures(lastMove Move) (bool, []Pos) {
	horizontalShift := lastMove.EndingPos.X - lastMove.InitialPos.X
	verticalShift := lastMove.EndingPos.Y - lastMove.InitialPos.Y
	var nextPieceF, nextPieceB Pos
	nextPieceF = Pos{lastMove.EndingPos.X + horizontalShift, lastMove.EndingPos.Y + verticalShift}
	nextPieceB = Pos{lastMove.InitialPos.X - horizontalShift, lastMove.InitialPos.Y - verticalShift}

	var captured = make([]Pos, 0)

	for board.IsValid(nextPieceF) {
		if board.BoardState[nextPieceF.X][nextPieceF.Y] == -lastMove.Player {

			captured = append(captured, nextPieceF)
		} else {
			break
		}
		nextPieceF = Pos{nextPieceF.X + horizontalShift, nextPieceF.Y + verticalShift}
	}
	for board.IsValid(nextPieceB) {
		if board.BoardState[nextPieceB.X][nextPieceB.Y] == -lastMove.Player {

			captured = append(captured, nextPieceB)
		} else {
			break
		}
		nextPieceB = Pos{nextPieceB.X + horizontalShift, nextPieceB.Y + verticalShift}
	}
	if len(captured) > 0 {
		return true, captured
	} else {
		return false, nil
	}
}

//create a new move if the positions are valid
func (board *Board) NewMove(player int, mFrom Pos, mTo Pos) Move {
	if player != board.BoardState[mFrom.X][mFrom.Y] {
		return Move{Invalid: errors.New("Players must move their own pieces")}

	}

	if board.BoardState[mTo.X][mTo.Y] != 0 {
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
				Player:     player,
				InitialPos: mFrom,
				EndingPos:  mTo,
				Invalid:    nil,
			}
			m.IsCapture, m.CappedPieces = board.CheckCaptures(m)
			return m
		} else {
			return Move{Invalid: errors.New("Invalid move")}
		}
	} else {
		//weak intersection
		if horizontalShift+verticalShift == 1 {
			m := Move{
				Player:     player,
				InitialPos: mFrom,
				EndingPos:  mTo,
				Invalid:    nil,
			}
			m.IsCapture, m.CappedPieces = board.CheckCaptures(m)
			return m

		} else {
			return Move{Invalid: errors.New("Invalid move")}
		}
	}

}

//applying a move
//returning the id of the player moving in the next turn and whether he can skip his move or not
//if the player id returned is -2 or 2, player -1 or 1 has won (respectively)
func (gameBoard *Board) ApplyMove(m Move) (int, bool) {
	gameBoard.BoardState[m.InitialPos.X][m.InitialPos.Y] = 0
	gameBoard.BoardState[m.EndingPos.X][m.EndingPos.Y] = m.Player
	//if no captures its the other players turn
	if !m.IsCapture {
		return -m.Player, false
	} else {
		for _, piece := range m.CappedPieces {
			gameBoard.BoardState[piece.X][piece.Y] = 0
		}
		//CHECK WIN
		white, black := gameBoard.CheckPieces()
		if white == 0 || black == 0 {
			return m.Player * 2, false
		}
		//MISSING CHECK IF NO LEGAL MOVES AVAILABLE
		return m.Player, true
	}

}

//print board to terminal
func (board *Board) PrintBoard() {
	fmt.Println("    0       1       2       3       4       5       6       7       8 ")
	for i, row := range board.BoardState {
		fmt.Printf("%d ", i)

		for j, piece := range row {
			if j < 8 {
				if piece == -1 {
					fmt.Printf(" %d  -- ", piece)

				} else {
					fmt.Printf("  %d  -- ", piece)

				}
			} else {
				if piece == -1 {
					fmt.Printf(" %d  ", piece)

				} else {
					fmt.Printf("  %d  ", piece)

				}
			}
		}
		if i < 4 {
			if i%2 == 0 {
				fmt.Println("\n    |   \\   |   /   |   \\   |   /   |   \\   |   /   |   \\   |   /   |")
				fmt.Println("\n    |     \\ | /     |     \\ | /     |     \\ | /     |     \\ | /     |")

			} else {
				fmt.Println("\n    |   /   |   \\   |   /   |   \\   |   /   |   \\   |   /   |   \\   |")
				fmt.Println("\n    | /     |     \\ | /     |     \\ | /     |     \\ | /     |     \\ |")

			}

		}

	}
}
