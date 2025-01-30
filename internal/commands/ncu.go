package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"strconv"
	"time"
)

func WaiterCustomize(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	selectedWaiter, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	newPriority, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return err
	}

	var w *objects.Waiter
	for _, w = range c.Location.Cafe().GetWaiters() {
		if w.ID == selectedWaiter {
			w.Name = req.Args[3]
			w.Priority = newPriority
			break
		}
	}

	c.Location.Broadcast("ncu", "-1", "0",
		req.Args[2],
		req.Args[3],
		req.Args[4],
	)

	// Stop waiter
	w.StopWorking()
	time.Sleep(100 * time.Millisecond)

	// Start it again
	w.IsWorking = true
	w.CurrentCounter = nil
	w.CurrentCustomer = nil
	go func() {
		for w.IsWorking {
			agents.IterateWaiter(c.Location, w)
		}
		w.CurrentCounter = nil
		w.Dish = -1
	}()

	return nil
}
