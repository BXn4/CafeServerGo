/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package marketplace

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
)

func init() {
	commands.RegisterCommand(requests.C2S_MARKETPLACE_JOBREFILL,
		commands.CommandConfig{
			Name:         "JobRefill",
			Identifier:   responses.S2C_MARKETPLACE_JOBREFILL,
			Description:  "Refilling available jobs",
			Args:         "{}",
			MinArgs:      0,
			MaxArgs:      0,
			FeatureLevel: 4,
		},
		MarketplaceJobRefillValidator,
		MarketplaceJobRefill,
		MarketPlaceJobRefillDBSaver,
	)
}

// mjo - MarketplaceJobRefill
func MarketplaceJobRefill(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	c.Player.AddGold(-balancing.BalancingConstants.JobRefillGold)
	c.Player.SetOpenJobs(5)

	c.Player.AddRefilledJobs()

	c.SendExtensionResponse(cm.Identifier, "-1", "0")
	return nil
}

func MarketplaceJobRefillValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Cant join the marketplace, because the player not yet reached the feature.", commands.NOT_DECLARED
	}

	// Check if player still has full or more than full jobs
	if c.Player.GetRefilledJobs() >= balancing.BalancingConstants.JobsPerDay {
		return "Player reached the jobs purchases / day!", commands.ERROR_JOB_REFILL_ONCE_A_DAY
	}

	// Check if player still has some jobs
	if c.Player.GetOpenJobs() != 0 {
		return "Player still have reaming jobs!", commands.NOT_DECLARED

	}

	// Check if player has enough gold
	if c.Player.GetGold() < balancing.BalancingConstants.JobRefillGold {
		return "Player not have enough money!", commands.NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
func MarketPlaceJobRefillDBSaver(c *client.Client) error {
	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())
	c.DB.UpdateRefilledJobs(c.Player.GetID(), c.Player.GetRefilledJobs())

	return nil
}
