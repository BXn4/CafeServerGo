package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
)

func init() {
	RegisterCommand(requests.C2S_MARKETPLACE_JOBREFILL,
		CommandConfig{
			Name:       "JobRefill",
			Identifier: responses.S2C_MARKETPLACE_JOBREFILL,
			MinArgs:    0,
			MaxArgs:    0,
		},
		MarketplaceJobRefillValidator,
		MarketplaceJobRefill,
	)
}

// mjo - MarketplaceJobRefill
func MarketplaceJobRefill(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.Player.AddGold(-balancing.BalancingConstants.JobRefillGold)
	c.Player.OpenJobs = 5

	c.Player.RefilledJobs = 1

	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())
	c.DB.UpdateRefilledJobs(c.Player.ID, c.Player.RefilledJobs)

	c.SendExtensionResponse("mjr", "-1", "0")
	return nil
}

func MarketplaceJobRefillValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	// Check if player still has full or more than full jobs
	if c.Player.RefilledJobs >= balancing.BalancingConstants.JobsPerDay {
		return "Player reached the jobs purchases / day!", ERROR_JOB_REFILL_ONCE_A_DAY
	}

	// Check if player still has some jobs
	if c.Player.OpenJobs != 0 {
		return "Player still have reaming jobs!", NOT_DECLARED

	}

	// Check if player has enough gold
	if c.Player.GetGold() < balancing.BalancingConstants.JobRefillGold {
		return "Player not have enough money!", NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", NO_ERROR
}
