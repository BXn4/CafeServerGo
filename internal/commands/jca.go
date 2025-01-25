package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// jca - JoinCafe
func JoinCafe(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Get id of cafe to join
	id, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	// Adds cafe to manager (loads it if not loaded)
	location := gm.AddLocation(id)

	// Send cafe joined
	c.SendExtensionResponse("jca", "-1", "0")

	// Leave cafe if already in one
	if c.Location != nil {
		c.Location.Leave(c.Player.ID)

		// Remove location if empty and owner is offline
		if c.Location.IsEmpty() && !gm.IsOnline(c.Location.Cafe().ID) && c.Location.Cafe().ID > 0 {

			gm.RemoveLocation(c.Location.Cafe().ID)
		}
	}

	// Join location
	location.Join(c.Player.ID, c.Writer)

	// Save location
	c.Location = location

	return nil
}
