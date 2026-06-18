package cafe

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/commands/editor"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_JOIN_CAFE,
		commands.CommandConfig{
			Name:        "JoinCafe",
			Identifier:  responses.S2C_JOIN_CAFE,
			Description: "Join an Cafe",
			Args:        "{}",
			MinArgs:     4,
			MaxArgs:     4,
		},
		JoinCafeValidator,
		JoinCafe,
		nil,
	)
}

// jca - JoinCafe
func JoinCafe(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	// Get id of cafe to join
	id, _ := strconv.Atoi(req.Args[3])

	// Adds cafe to manager (loads it if not loaded)
	location, err := gm.AddLocation(id)

	if err != nil {
		return fmt.Errorf("Failed to load location %d: %v", id, err)
	}

	// Leave cafe if already in one
	if c.Location != nil {
		c.Location.Leave(c.Player.GetID())
	}

	// Send cafe joined
	c.SendExtensionResponse(cm.Identifier, "-1", "0")

	// Save location
	c.Location = location

	// Send fridge info (ifr)
	SendFridgeInventory(req, c, gm)

	editor.SendFurnitureInventory(req, c, gm)

	// Join location
	c.Location.Join(c.Player.GetID(), c.ResponseQueue)

	return nil
}

func JoinCafeValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	_, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR // CAFE_NOT_EXIST
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
