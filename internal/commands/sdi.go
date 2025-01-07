package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// sdi - C2S_SHOP_DELETE_ITEM
func SellIngrediement(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {
	/*ingrediementID := req.Args[1]
	ingrediementAmount := req.Args[2]*/

	c.SendExtensionResponse("pin", "-1")

	return nil
}
