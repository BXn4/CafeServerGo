package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// mjm - JoinMarketplace
func JoinMarketplace(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us MJM while in editor.
	if c.Location.Cafe().InEditorMode() {
		return nil
	}

	// Gets cafe location
	location := gm.AddLocation(-1)

	// Send cafe joined
	c.SendExtensionResponse("mjm", "-1", "0")

	// Leave current cafe if there is one
	if c.Location != nil {
		c.Location.Leave(c.Player.ID)

		// Remove location if its empty and owner not online
		if c.Location.IsEmpty() && !gm.IsOnline(c.Location.Cafe().GetID()) && c.Location.Cafe().GetID() > 0 {
			gm.RemoveLocation(c.Location.Cafe().GetID())
		}
	}

	// Join cafe
	location.Join(c.Player.ID, c.ResponseQueue)

	// Save location
	c.Location = location

	return nil
}
