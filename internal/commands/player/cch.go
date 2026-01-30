package player

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/commands/cmdlet"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_CHAT,
		commands.CommandConfig{
			Name:        "Chat",
			Identifier:  responses.S2C_CAFE_CHAT,
			Description: "Broadcasting chat messages",
			Args:        "{FromPlayerID} {Message}",
			MinArgs:     3,
			MaxArgs:     3,
		},
		ChatMessageValidator,
		SendChatMessage,
		nil,
	)
}

// cch - C2S_CAFE_CHAT
func SendChatMessage(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	message := req.Args[2]

	if strings.HasPrefix(message, "/") {
		cmdlet.HandleCmdlets(c, gm, message)
		return nil
	}

	println(c.Player.GetPos().X)

	c.Location.Broadcast(cm.Identifier, "-1", "0", strconv.Itoa(c.Player.GetID()), message)

	return nil
}

func ChatMessageValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CCH while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", commands.LOCATION_NOT_RUNNING
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
