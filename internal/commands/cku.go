package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// cku - KickPlayer
func KickPlayer(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Check if client is the location owner
	if c.Location.Cafe().GetID() != c.Player.ID {
		return nil
	}

	id, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	// If not at location return
	if !c.Location.AtLocation(id) {
		return nil
	}

	// Send kicked out message
	c.Location.Send(id, "cku", "-1", req.Args[2])

	// Get clien of the kicked out user and leave room
	item, err := gm.GetClient(id)
	if err != nil {
		return err
	}
	kickedClient := item.(*client.Client)
	kickedClient.Location.Leave(id)

	// Get owned location of the kicked out user
	location := gm.AddLocation(id)

	// Join kicked out user to its owned cafe
	location.Join(kickedClient.Player.ID, kickedClient.ResponseQueue)

	// Set new location
	kickedClient.Location = location

	return nil
}
