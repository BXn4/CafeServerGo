package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/shop"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"strconv"
	"strings"
)

func init() {
	RegisterCommand(requests.C2S_SHOP_AVAILIBILITY,
		CommandConfig{
			Name:       "ShopAvailibility",
			Identifier: responses.S2C_SHOP_AVAILIBILITY,
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		ShopAvailibility,
	)
}

// sga - C2S_SHOP_AVAILIBILITY
func ShopAvailibility(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	var args []string
	// checks from the current time
	if shop.IsShopUnavailable() {
		for _, id := range shop.GetUnavailableIngredients() {
			args = append(args, strconv.Itoa(id))
		}
	}

	c.SendExtensionResponse("sga", "-1", "0", strings.Join(args, "#"))

	return nil
}
