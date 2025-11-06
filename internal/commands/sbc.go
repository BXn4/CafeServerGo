package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
)

// sbc - SendBalancingConstant
func SendBalancingConstant(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	c.SendExtensionResponse("sbc", "1", "0", balancing.BalancingConstants.AsResponse())

	return nil
}
