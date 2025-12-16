package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
)

// min level: 7

func init() {
	RegisterCommand(requests.C2S_GIFT_PLAYERGIFTS,
		CommandConfig{
			Name:       "PlayerGifts",
			Identifier: responses.S2C_GIFT_PLAYERGIFTS,
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		SendPlayerGifts,
	)
}

func SendPlayerGifts(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendExtensionResponse("gmg", "-1", "0", c.Player.Gifts.StringWithIndex())
	return nil
}
