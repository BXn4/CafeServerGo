package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
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

	// Leave cafe if already in one
	if c.Location != nil {
		c.Location.Leave(c.Player.ID)
	}

	// Send cafe joined
	c.SendExtensionResponse("jca", "-1", "0")

	// Save location
	c.Location = location

	// Send fridge info (ifr)
	err = SendFridgeInventory(req, c, gm)
	if err != nil {
		return fmt.Errorf("\nifr request: %v", err)
	}

	// Join location
	c.Location.Join(c.Player.ID, c.ResponseQueue)

	return nil
}
