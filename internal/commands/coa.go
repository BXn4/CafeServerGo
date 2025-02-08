package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// coa - CoopActiveList
func CoopActiveList(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	c.SendExtensionResponse("coa", "-1", "0", "2101")
	return nil
}
