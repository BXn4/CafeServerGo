package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// sga - C2S_SHOP_AVAILIBILITY
func ShopAvailibility(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	c.SendExtensionResponse("sga", "-1", "0", "1314#1315")
	return nil
}
