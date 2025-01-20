package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// pin
func FireWaiter(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	index, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}
	index -= 1

	// TODO: Make waiter leave

	c.Location.Cafe().Waiters = append(c.Location.Cafe().Waiters[:index], c.Location.Cafe().Waiters[index+1:]...)

	c.SendExtensionResponse("nfi", "0", "0", req.Args[2])

	return nil
}
