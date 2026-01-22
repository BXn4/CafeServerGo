package player

import (
	"cafego/internal/client"
	"cafego/internal/types/responses"
	"strconv"
)

// asy - ASSETS_SYNCHRONIZE
// Updates force the player cash, gold in the game visually.
// Its used in the payments, but we can use it to update the cash, gold values force.

func AssetsSynchronize(c *client.Client) error {
	c.SendExtensionResponse(responses.S2C_ASSETS_SYNCHRONIZE, "-1", "0", strconv.Itoa(c.Player.GetCash()), strconv.Itoa(c.Player.GetGold()))
	return nil
}
