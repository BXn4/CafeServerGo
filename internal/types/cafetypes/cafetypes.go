package cafetypes

type CafeObjectKind int

const (
	STOVE   CafeObjectKind = 0
	COUNTER                = 1
	CHAIR                  = 2
	TABLE                  = 3
	VENDING                = 4
	OTHER                  = 5
)

type CafeObjectRotation int

const (
	Up CafeObjectRotation = iota
	Left
	Down
	Right
)
