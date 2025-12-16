package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"time"
)

func init() {
	RegisterCommand(requests.C2S_PING,
		CommandConfig{
			Name:       "ServerPin",
			Identifier: responses.S2C_PING,
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		SendPing,
	)
}

// pin
func SendPing(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.TimeoutStamp = time.Now()
	c.SendExtensionResponse("pin", "-1")
	return nil
}
