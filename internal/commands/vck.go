package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/versions"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_VERSION_CHECK,
		CommandConfig{
			Name:       "VersionCheck",
			Identifier: responses.S2C_VERSION_CHECK,
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		VersionCheck,
	)
}

// vck - VersionCheck
func VersionCheck(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	version := strconv.Itoa(versions.GetGameVersion())
	c.SendExtensionResponse("vck", "1", "0", version)
	return nil
}
