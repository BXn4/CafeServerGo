package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

// abr - AllowFriendRequests
func AllowEmails(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	c.Player.AllowEmails = utils.If(req.Args[2] == "1", true, false)

	c.Location.Broadcast("abr", "-1", req.Args[2], strconv.Itoa(c.Player.ID))
	return nil
}
