package commands

import (
	"cafego/internal/client"
	"cafego/internal/commands/admin"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
	"strings"
)

// cch - C2S_CAFE_CHAT
func SendChatMessage(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us CCH while in editor.
	if c.Location.Cafe().InEditorMode() {
		return nil
	}

	message := req.Args[2]

	if strings.HasPrefix(message, "/") {
		admin.HandleAdminCommands(req, c, gm, message)
		return nil
	}

	c.Location.Broadcast("cch", "-1", "0", strconv.Itoa(c.Player.ID), message)

	return nil
}
