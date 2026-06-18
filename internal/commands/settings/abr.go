package settings

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_ALLOW_BUDDY_REQUESTS,
		commands.CommandConfig{
			Name:        "AllowBuddyRequests",
			Identifier:  responses.S2C_ALLOW_BUDDY_REQUESTS,
			Description: "Enable/Disable friend requests.",
			Args:        "{0/1}",
			MinArgs:     3,
			MaxArgs:     3,
			IsBool:      true,
		},
		AllowFriendRequestsValidator,
		AllowFriendRequests,
		AllowFriendRequestsDBSaver,
	)
}

// abr - AllowFriendRequests
func AllowFriendRequests(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	c.Player.SetAllowFriendRequests(utils.If(req.Args[2] == "1", true, false))
	c.Location.Broadcast(cm.Identifier, "-1", req.Args[2], strconv.Itoa(c.Player.GetID()))
	return nil
}

func AllowFriendRequestsValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if cm.IsBool {
		if req.Args[2] != "0" && req.Args[2] != "1" {
			return fmt.Sprintf("Invalid args for boolean: %v", req.Args[2]), commands.INVALID_ARGS
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func AllowFriendRequestsDBSaver(c *client.Client) error {
	c.DB.UpdateAllowFriendRequests(c.Player.GetID(), c.Player.GetAllowFriendRequests())

	return nil
}
