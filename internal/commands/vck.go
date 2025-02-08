package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// vck - VersionCheck
func VersionCheck(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendExtensionResponse("vck", "1", "0", "1603")
	return nil
}
