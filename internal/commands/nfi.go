package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// pin
func FireWaiter(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	id, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	c.Location.Cafe().RemoveWaiter(id)

	c.SendExtensionResponse("nfi", "0", "0", req.Args[2])

	return nil
}
