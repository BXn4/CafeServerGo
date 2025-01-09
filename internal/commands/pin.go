package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// pin
func SendPing(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	c.SendExtensionResponse("pin", "-1")

	return nil
}
