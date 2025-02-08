package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/waiter"
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

	var w *waiter.Waiter
	for _, w = range c.Location.Cafe().GetWaiters() {
		if w.GetID() == selectedWaiter {
			a := w.GetAvatar()
			a.Name = req.Args[3]
			w.SetAvatar(a)
			w.SetPriority(newPriority)
			break
		}
	}

	c.Location.Broadcast("ncu", "-1", "0",
		req.Args[2], // waiter id
		req.Args[3], // new waiter name
		req.Args[4], // new waiter priority
	)

	// Stop waiter
	println("DEBUG")
	w.StopWorking()
	time.Sleep(100 * time.Millisecond)

	c.Location.Broadcast("nav", "-1", "0", w.SpawnString())

	return nil
}
