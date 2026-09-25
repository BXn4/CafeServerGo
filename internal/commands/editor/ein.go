/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package editor

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
)

// ein - SendFurnitureInventory
func SendFurnitureInventory(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	inventory := c.Location.Cafe().GetFurnitureInventory().String()

	c.SendExtensionResponse(responses.S2C_EDITOR_INVENTORY, "1", "0",
		inventory,
	)
	return nil
}
