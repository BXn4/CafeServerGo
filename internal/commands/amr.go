package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

// abr - AllowEmails
func AllowEmails(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	c.Player.AllowEmails = utils.If(req.Args[2] == "1", true, false)

	c.SendExtensionResponse("amr", "-1", req.Args[2], strconv.Itoa(c.Player.ID))

	c.DB.UpdateAllowEmails(c.Player.ID, c.Player.AllowEmails)
	return nil
}
