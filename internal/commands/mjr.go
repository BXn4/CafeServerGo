package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
)

// vck - VersionCheck
func MarketplaceJobRefill(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Check if player still has full or more than full jobs
	if c.Player.OpenJobs >= balancing.BalancingConstants.JobsPerDay {
		c.SendExtensionResponse("mjr", "-1", "4")
		return nil
	}

	// Check if player still has some jobs
	if c.Player.OpenJobs != 0 {
		c.SendExtensionResponse("mjr", "-1", "1")
		return nil
	}

	// Check if player has enough gold
	if c.Player.GetGold() < balancing.BalancingConstants.JobRefillGold {
		return nil
	}

	c.Player.AddGold(-balancing.BalancingConstants.JobRefillGold)
	c.Player.OpenJobs = balancing.BalancingConstants.JobsPerDay

	c.SendExtensionResponse("mjr", "-1", "0")
	return nil
}
