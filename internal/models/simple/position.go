package simple

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

// Removed the Scan, and the Value
// Because PlayerStart always starts from the DOOR,
// and we can it by the door, so dont need to save it
