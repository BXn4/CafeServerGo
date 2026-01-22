package gifts

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
)

func SendPlayerGifts(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	gifts := c.Player.GetGifts()
	c.SendExtensionResponse(responses.S2C_GIFT_PLAYERGIFTS, "-1", "0", gifts.StringWithIndex())
	return nil
}
