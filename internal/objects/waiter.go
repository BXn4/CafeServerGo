package objects

type Action int

const (
	WAITER_INSERT Action = iota
	WAITER_MOVE_TO_COUNTER
	WAITER_FEED
	WAITER_CLEAN
)

type Waiter struct {
	ID       int
	Name     string
	Priority int
	Avatar   Avatar
	Pos      []int
	Counter  CafeObject
	Dish     int
	Action   Action
	Task     string //TODO: Change it
}
