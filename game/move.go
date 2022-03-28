package game

//move struct
type Move struct {
	Player     int
	InitialPos Pos
	EndingPos  Pos
	IsCapture  bool
	CappedFwd  []Pos
	CappedBwd  []Pos

	Invalid    error
	ReCaptures []Move
}
