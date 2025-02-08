package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

func SendPlayerGifts(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendExtensionResponse("gmg", "-1", "0", c.Player.Gifts.StringWithIndex())
	return nil
}
