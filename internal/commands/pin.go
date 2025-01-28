package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"time"
)

// pin
func SendPing(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.TimeoutStamp = time.Now()
	c.SendExtensionResponse("pin", "-1")
	return nil
}
