package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// lmi - SendMasteryInfo
func SendMasteryInfo(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendExtensionResponse("lmi", "-1", "0", c.Player.BuildMastery())
	return nil
}
