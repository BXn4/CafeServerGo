package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
	"strings"
)

// sga - C2S_SHOP_AVAILIBILITY
func ShopAvailibility(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	var args []string
	// checks from the current time
	if gm.IsShopUnavailable() {
		for _, id := range gm.GetUnavailableIngredients() {
			args = append(args, strconv.Itoa(id))
		}
	}

	c.SendExtensionResponse("sga", "-1", "0", strings.Join(args, "#"))

	return nil
}
