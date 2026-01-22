package player

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// lmi - SendMasteryInfo
func SendMasteryInfo(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	mastery := c.Player.GetMastery()
	c.SendExtensionResponse("lmi", "-1", "0", mastery)
	return nil
}
