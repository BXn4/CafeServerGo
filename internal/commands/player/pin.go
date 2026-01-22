package player

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"time"
)

func init() {
	commands.RegisterCommand(requests.C2S_PING,
		commands.CommandConfig{
			Name:        "ServerPin",
			Identifier:  responses.S2C_PING,
			Description: "Ping",
			Args:        "{}",
			MinArgs:     0,
			MaxArgs:     0,
		},
		nil,
		SendPing,
		nil,
	)
}

// pin
func SendPing(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	c.TimeoutStamp = time.Now()
	c.SendExtensionResponse(responses.S2C_PING, "-1")
	return nil
}
