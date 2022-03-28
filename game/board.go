package game

import (
	"bufio"
	"fmt"
)

//general board state struct
type Board struct {
	BoardState [5][9]int
	Rows       int
	Cols       int
}

func requestPosition(inputBuffer bufio.Reader) Pos {
	return Pos{}
}

func requestDirection(inputBuffer bufio.Reader) rune {
	return 'F' //or 'B' for declaring forward or backwards captures
}

//copy existing game
//func CopyGame(existingGame *Game) *Game

//create a default starting board
func NewBoard() *Board {
	board := Board{
		BoardState: [5][9]int{
			{-1, -1, -1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1, -1, -1},
			{-1, 1, -1, 1, 0, -1, 1, -1, 1},
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
func (board *Board) CheckCaptures(lastMove *Move) {
	horizontalShift := lastMove.EndingPos.X - lastMove.InitialPos.X
	verticalShift := lastMove.EndingPos.Y - lastMove.InitialPos.Y
	var nextPieceF, nextPieceB Pos
	nextPieceF = Pos{lastMove.EndingPos.X + horizontalShift, lastMove.EndingPos.Y + verticalShift}
	nextPieceB = Pos{lastMove.InitialPos.X - horizontalShift, lastMove.InitialPos.Y - verticalShift}

	var frontCaptured = make([]Pos, 0)
	var backCaptured = make([]Pos, 0)

	for board.IsValid(nextPieceF) {
		if board.BoardState[nextPieceF.X][nextPieceF.Y] == -lastMove.Player {

			frontCaptured = append(frontCaptured, nextPieceF)
		} else {
			break
		}
		nextPieceF = Pos{nextPieceF.X + horizontalShift, nextPieceF.Y + verticalShift}
	}
	for board.IsValid(nextPieceB) {
		if board.BoardState[nextPieceB.X][nextPieceB.Y] == -lastMove.Player {

			backCaptured = append(backCaptured, nextPieceB)
		} else {
			break
		}
		nextPieceB = Pos{nextPieceB.X - horizontalShift, nextPieceB.Y - verticalShift}
	}
	if len(frontCaptured) > 0 || len(backCaptured) > 0 {
		lastMove.IsCapture = true

	} else {
		lastMove.IsCapture = false
	}
	lastMove.CappedFwd = frontCaptured
	lastMove.CappedBwd = backCaptured
}

//not yet necessary, might delete
func (gameBoard *Board) CopyBoard() *Board {
	newB := Board{
		BoardState: gameBoard.BoardState,
		Rows:       gameBoard.Rows,
		Cols:       gameBoard.Cols,
	}
	return &newB
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

//print board to terminal
func (board *Board) ToString() string {
	var str string
	str += "    0       1       2       3       4       5       6       7       8 \n"
	for i, row := range board.BoardState {
		str += fmt.Sprintf("%d ", i)

		for j, piece := range row {
			if j < 8 {
				if piece == -1 {
					str += fmt.Sprintf(" %d  -- ", piece)

				} else {
					str += fmt.Sprintf("  %d  -- ", piece)

				}
			} else {
				if piece == -1 {
					str += fmt.Sprintf(" %d  ", piece)

				} else {
					str += fmt.Sprintf("  %d  ", piece)

				}
			}
		}
		if i < 4 {
			if i%2 == 0 {
				str += fmt.Sprintf("\n    |   \\   |   /   |   \\   |   /   |   \\   |   /   |   \\   |   /   |\n")
				str += fmt.Sprintf("\n    |     \\ | /     |     \\ | /     |     \\ | /     |     \\ | /     |\n")

			} else {
				str += fmt.Sprintf("\n    |   /   |   \\   |   /   |   \\   |   /   |   \\   |   /   |   \\   |\n")
				str += fmt.Sprintf("\n    | /     |     \\ | /     |     \\ | /     |     \\ | /     |     \\ |\n")

			}

		}

	}
	return str
}
