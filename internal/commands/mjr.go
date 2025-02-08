package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// vck - VersionCheck
func MarketplaceJobRefill(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Check if player still has full or more than full jobs
	if c.Player.OpenJobs >= 5 {
		c.SendExtensionResponse("mjr", "-1", "4")
		return nil
	}

	// Check if player still has some jobs
	if c.Player.OpenJobs != 0 {
		c.SendExtensionResponse("mjr", "-1", "1")
		return nil
	}

	// Check if player has enough gold
	if c.Player.GetGold() < 1 {
		return nil
	}

	c.Player.AddGold(-1)
	c.Player.OpenJobs = 5

	c.SendExtensionResponse("mjr", "-1", "0")
	return nil
}
