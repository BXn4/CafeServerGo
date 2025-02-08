package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// ifr - SendFridgeInventory
func SendFridgeInventory(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	fridge := c.Location.Cafe().GetFridgeInventory()
	fridgeCap := c.Location.Cafe().GetFridgeMaxCapacity()

	c.SendExtensionResponse("ifr", "1", "0",
		strconv.Itoa(fridgeCap),
		fridge.String(),
	)
	return nil
}
