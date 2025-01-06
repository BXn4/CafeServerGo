package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// lmi - SendMasteryInfo
func SendMasteryInfo(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

	println("PLAYER MASTERY: ", c.Player.Mastery)

	c.SendExtensionResponse("lmi", "-1", "0", c.Player.Mastery)

	return nil
}
