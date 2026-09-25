/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package gifts

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
)

func init() {
	commands.RegisterCommand(requests.C2S_GIFT_PLAYERGIFTS,
		commands.CommandConfig{
			Name:         "Player gifts",
			Description:  "List of player gifts",
			Args:         "{gifts}",
			Identifier:   responses.S2C_GIFT_PLAYERGIFTS,
			MinArgs:      0,
			MaxArgs:      0,
			FeatureLevel: 7,
		},
		nil,
		SendPlayerGifts,
		nil,
	)
}

func SendPlayerGifts(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	c.SendExtensionResponse(responses.S2C_GIFT_PLAYERGIFTS, "-1", "0", c.Player.GetGiftStringWithIndex())
	return nil
}
