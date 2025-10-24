package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// ein - SendFurnitureInventory
func SendFurnitureInventory(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	inventory := c.Location.Cafe().FurnitureInventory.String()

	c.SendExtensionResponse("ein", "1", "0",
		inventory,
	)
	return nil
}
