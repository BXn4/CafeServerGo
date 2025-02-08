package commands

import (
	"cafego/internal/client"
	"cafego/internal/commands/cmdlet"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
	"strings"
)

// cch - C2S_CAFE_CHAT
func SendChatMessage(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us CCH while in editor.
	if !c.Location.IsRunning() {
		return nil
	}

	message := req.Args[2]

	if strings.HasPrefix(message, "/") {
		cmdlet.HandleCmdlets(c, gm, message)
		return nil
	}

	c.Location.Broadcast("cch", "-1", "0", strconv.Itoa(c.Player.ID), message)

	return nil
}
