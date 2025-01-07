package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// sga - C2S_SHOP_AVAILIBILITY
func ShopAvailibility(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {
	// TODO: Set what ingredients are not available.

	c.SendExtensionResponse("sga", "-1", "0", "#")

	return nil
}
