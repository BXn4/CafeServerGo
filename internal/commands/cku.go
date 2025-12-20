package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_KICK_USER,
		CommandConfig{
			Name:       "KickPlayerFromRoom",
			Identifier: responses.S2C_KICK_USER,
			MinArgs:    3,
			MaxArgs:    3,
		},
		KickPlayerValidator,
		KickPlayer,
	)
}

// cku - KickPlayer
func KickPlayer(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	id, _ := strconv.Atoi(req.Args[2])

	// Send kicked out message
	c.Location.Send(id, "cku", "-1", req.Args[2])

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
	location.Join(kickedClient.Player.ID, kickedClient.ResponseQueue)

	// Set new location
	kickedClient.Location = location

	return nil
}

func KickPlayerValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	id, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	// Check if client is the location owner
	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the location owner!", NOT_DECLARED
	}

	// If not at location return
	if !c.Location.AtLocation(c.Location.Cafe().ID) {
		return "Not at the location!", NOT_DECLARED
	}

	// Get client of the kicked out user and leave room
	_, err = gm.GetClient(id)
	if err != nil {
		return "Could not get the client!", NOT_DECLARED

	}

	return "Command ran without any errors.", NO_ERROR
}
