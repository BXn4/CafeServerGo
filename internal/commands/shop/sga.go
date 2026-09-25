/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package shop

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/shop"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"strconv"
	"strings"
)

func init() {
	commands.RegisterCommand(requests.C2S_SHOP_AVAILIBILITY,
		commands.CommandConfig{
			Name:        "ShopAvailibility",
			Identifier:  responses.S2C_SHOP_AVAILIBILITY,
			Description: "Shop availibility",
			Args:        "{ingredients}#",
			MinArgs:     0,
			MaxArgs:     0,
		},
		nil,
		ShopAvailibility,
		nil,
	)
}

// sga - C2S_SHOP_AVAILIBILITY
func ShopAvailibility(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	var args []string
	// checks from the current time
	if shop.IsShopUnavailable() {
		for _, id := range shop.GetUnavailableIngredients() {
			args = append(args, strconv.Itoa(id))
		}
	}

	c.SendExtensionResponse(responses.S2C_SHOP_AVAILIBILITY, "-1", "0", strings.Join(args, "#"))

	return nil
}
