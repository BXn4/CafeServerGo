package player

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/versions"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_VERSION_CHECK,
		commands.CommandConfig{
			Name:        "VersionCheck",
			Identifier:  responses.S2C_VERSION_CHECK,
			Description: "Check the client / server game version",
			Args:        "{version}",
			MinArgs:     0,
			MaxArgs:     0,
		},
		nil,
		VersionCheck,
		nil,
	)
}

// vck - VersionCheck
func VersionCheck(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	version := strconv.Itoa(versions.GetGameVersion())
	c.SendExtensionResponse(cm.Identifier, "1", "0", version)
	return nil
}
