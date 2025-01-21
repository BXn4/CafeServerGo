package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
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

	for _, w := range c.Location.Cafe().Waiters {
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
	return nil
}
