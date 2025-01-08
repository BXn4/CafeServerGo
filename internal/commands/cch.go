package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// cch - C2S_CAFE_CHAT
func SendChatMessage(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

	message := req.Args[2]

	c.Cafe.Broadcast("cch", "-1", "0", strconv.Itoa(c.Player.ID), message)

	return nil
}
