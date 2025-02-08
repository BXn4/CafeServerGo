package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

func RemoveGift(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	slot, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	c.Player.Gifts.RemoveGift(slot)

	c.SendExtensionResponse("gmg", "-1", "0", req.Args[2])

	SendPlayerGifts(req, c, gm)
	return nil
}
