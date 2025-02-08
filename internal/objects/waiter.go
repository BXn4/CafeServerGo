package objects

import "fmt"

const (
	WAITER_INSERT Action = iota
	WAITER_MOVE_TO_COUNTER
	WAITER_FEED
	WAITER_CLEAN
)

type Priority int

const (
	CLEANING = 0
	BOTH     = 50
	SERVING  = 100
)

type Action int

const (
	INSERT          Action = 0
	MOVE_TO_COUNTER        = 5
	FEED                   = 6
	CLEAN                  = 7
)

type Waiter struct {
	ID              int
	Name            string
	Priority        int
	Avatar          Avatar
	Pos             [2]int
	Dish            int
	Action          Action
	CurrentCounter  *CafeObject
	CurrentCustomer *Customer
	IsWorking       bool
}

func (w *Waiter) StopWorking() {
	w.IsWorking = false
	if w.CurrentCustomer != nil {
		w.CurrentCustomer.SetAssignedWaiter(-1)
	}
}

func (w *Waiter) String() string {
	return fmt.Sprintf("%v+%v+%v", w.Name, w.Avatar.Apperance(), w.Priority)
}
