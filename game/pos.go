package game

//board position struct
type Pos struct {
	X int
	Y int
}

func equals(p1 Pos, p2 Pos) bool {
	if p1.X == p2.X && p1.Y == p2.Y {
		return true
	} else {
		return false
	}
}
