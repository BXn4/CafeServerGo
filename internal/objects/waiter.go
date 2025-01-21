package objects

import (
	"cafego/internal/types/daos"
	"encoding/json"
	"fmt"
)

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
	Pos             []int
	Dish            int
	Action          Action
	CurrentCounter  *CafeObject
	CurrentCustomer *Customer
	IsWorking       bool
}

func (w *Waiter) StopWorking() {
	w.IsWorking = false
	if w.CurrentCustomer != nil {
		w.CurrentCustomer.AssignedWaiter = -1
	}
}

func (w *Waiter) JSON() string {

	dao := daos.WaiterDAO{
		ID:       w.ID,
		Name:     w.Name,
		Priority: w.Priority,
		Avatar:   w.Avatar.String(),
	}

	b, err := json.Marshal(dao)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
