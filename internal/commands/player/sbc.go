package player

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
)

// sbc - SendBalancingConstant
func SendBalancingConstant(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	c.SendExtensionResponse(responses.S2C_SEND_BLANCING_CONSTANTS, "1", "0", balancing.BalancingConstants.AsResponse())

	return nil
}
