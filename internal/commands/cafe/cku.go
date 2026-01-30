package cafe

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_KICK_USER,
		commands.CommandConfig{
			Name:        "KickPlayerFromRoom",
			Identifier:  responses.S2C_KICK_USER,
			Description: "Kicks player from room",
			Args:        "{id} {id}",
			MinArgs:     3,
			MaxArgs:     3,
		},
		KickPlayerValidator,
		KickPlayer,
		nil,
	)
}

// cku - KickPlayer
func KickPlayer(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	id, _ := strconv.Atoi(req.Args[2])

	// Send kicked out message
	c.Location.Send(id, cm.Identifier, "-1", req.Args[2])

	// Get clien of the kicked out user and leave room
	item, _ := gm.GetClient(id)

	kickedClient := item.(*client.Client)
	kickedClient.Location.Leave(id)

	// Get owned location of the kicked out user
	location, err := gm.AddLocation(id)
	if err != nil {
		return fmt.Errorf("Failed to load location for kicked user %d: %v", id, err)
	}

	// Join kicked out user to its owned cafe
	location.Join(kickedClient.Player.GetID(), kickedClient.ResponseQueue)

	// Set new location
	kickedClient.Location = location

	return nil
}

func KickPlayerValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	id, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	// Check if client is the location owner
	if c.Location.Cafe().GetOwnerID() != c.Player.GetID() {
		return "Not the location owner!", commands.NOT_DECLARED
	}

	// If not at location return
	if !c.Location.AtLocation(c.Location.Cafe().GetID()) {
		return "Not at the location!", commands.NOT_DECLARED
	}

	// Get client of the kicked out user and leave room
	_, err = gm.GetClient(id)
	if err != nil {
		return "Could not get the client!", commands.NOT_DECLARED

	}

	return "Command ran without any errors.", commands.NO_ERROR
}
