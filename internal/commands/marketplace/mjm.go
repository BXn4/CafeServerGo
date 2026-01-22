package marketplace

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
)

func init() {
	commands.RegisterCommand(requests.C2S_MARKETPLACE_JOIN,
		commands.CommandConfig{
			Name:         "JoinMarkeplace",
			Identifier:   responses.S2C_MARKETPLACE_JOIN,
			Description:  "Joining to the marketplace",
			Args:         "{}",
			MinArgs:      0,
			MaxArgs:      0,
			FeatureLevel: 4,
		},
		JoinMarketplaceValidator,
		JoinMarketplace,
		nil,
	)
}

// min level 4

// mjm - JoinMarketplace
func JoinMarketplace(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {

	// Gets cafe location
	location, err := gm.AddLocation(-1)
	if err != nil {
		return fmt.Errorf("Failed to load marketplace: %v", err)
	}

	// Send cafe joined
	c.SendExtensionResponse(cm.Identifier, "-1", "0")

	// Leave current cafe if there is one
	if c.Location != nil {
		c.Location.Leave(c.Player.GetID())
	}

	// Join cafe
	location.Join(c.Player.GetID(), c.ResponseQueue)

	// Save location
	c.Location = location

	return nil
}

func JoinMarketplaceValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
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

	// Dont allow players to modify the packet and sending us MJM while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", commands.LOCATION_NOT_RUNNING
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
