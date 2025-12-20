package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
)

func init() {
	RegisterCommand(requests.C2S_MARKETPLACE_JOIN,
		CommandConfig{
			Name:       "JoinMarkeplace",
			Identifier: responses.S2C_MARKETPLACE_JOIN,
			MinArgs:    0,
			MaxArgs:    0,
		},
		JoinMarketplaceValidator,
		JoinMarketplace,
	)
}

// min level 4

// mjm - JoinMarketplace
func JoinMarketplace(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Gets cafe location
	location, err := gm.AddLocation(-1)
	if err != nil {
		return fmt.Errorf("Failed to load marketplace: %v", err)
	}

	// Send cafe joined
	c.SendExtensionResponse("mjm", "-1", "0")

	// Leave current cafe if there is one
	if c.Location != nil {
		c.Location.Leave(c.Player.ID)
	}

	// Join cafe
	location.Join(c.Player.ID, c.ResponseQueue)

	// Save location
	c.Location = location

	return nil
}

func JoinMarketplaceValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if c.Player.GetLevel() < 4 {
		return "Cant join the marketplace, because the player not yet reached the feature.", NOT_DECLARED
	}

	// Dont allow players to modify the packet and sending us MJM while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", LOCATION_NOT_RUNNING
	}

	return "Command ran without any errors.", NO_ERROR
}
